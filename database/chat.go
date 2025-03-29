package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"html/template"
	"net/url"
	"strings"
	"time"
)

type Message struct {
	SenderName string
	SendDate   time.Time
	Text       template.HTML
}

func (m *Message) GetTimeSend(a *Account) string {
	return m.SendDate.In(a.TimeZone).Format(loc.TimeFormatString)
}

type ChatRoom struct {
	Name       string
	Member     []string
	User       string
	NewMessage bool
}

func (c *ChatRoom) GetLink() template.URL {
	return template.URL("/chat/" + url.PathEscape(c.Name) + "/" + url.PathEscape(c.User))
}

func (c *ChatRoom) GetMemberList() string {
	return strings.Join(c.Member, ", ")
}

func LoadLastMessages(amount int, timeStamp time.Time, roomID string, accountName string) ([]Message, error) {
	err := postgresDB.QueryRow(`SELECT account_name from chat_rooms_to_account WHERE account_name = $1 AND room_id = $2`, accountName, roomID).Scan(&accountName)
	if err != nil {
		return nil, err
	}
	result, err := postgresDB.Query(`SELECT sender, send_time, message FROM chat_messages 
WHERE send_time < $1 AND room_id = $2 ORDER BY send_time DESC LIMIT $3`, timeStamp, roomID, amount)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]Message, 0)
	msg := Message{}
	for result.Next() {
		err = result.Scan(&msg.SenderName, &msg.SendDate, &msg.Text)
		if err != nil {
			return nil, err
		}
		arr = append(arr, msg)
	}
	return arr, nil
}

func InsertMessage(msg *Message, roomID string) error {
	_, err := postgresDB.Exec(`INSERT INTO chat_messages (room_id, sender, message, send_time) VALUES ($1, $2, $3, $4)`,
		roomID, msg.SenderName, msg.Text, msg.SendDate)
	return err

}

func CreateChatRoom(roomID string, member []string) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	var name string
	err = tx.QueryRow(`SELECT room_id FROM chat_rooms WHERE member = 
	ARRAY(SELECT name FROM account WHERE name = ANY($1) AND name <> $2 AND blocked = false ORDER BY name)`,
		pq.Array(member), loc.AdministrationAccountName).Scan(&name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err == nil {
		return DoubleChatRoomEntry
	}

	err = tx.QueryRow(`SELECT room_id FROM chat_rooms WHERE room_id = $1`, roomID).Scan(&name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err == nil {
		return ChatRoomNameTaken
	}

	_, err = tx.Exec(`INSERT INTO chat_rooms (room_id, member, created) VALUES ($1, 
	ARRAY(SELECT name FROM account WHERE name = ANY($2) AND name <> $3 AND blocked = false ORDER BY name), $4)`,
		roomID, pq.Array(member), loc.AdministrationAccountName, time.Now().UTC())
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO chat_rooms_to_account (room_id, account_name, new_message) 
SELECT $1 AS room_id, name, false AS new_message FROM account WHERE name = ANY($2)`, roomID, pq.Array(member))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func QueryForRoomIdAndUser(roomID string, accountName string, ownerName string) error {
	return postgresDB.QueryRow(`SELECT room_id FROM chat_rooms_linked 
               WHERE room_id = $1 AND account_name = $2 AND owner_name = $3;`,
		roomID, accountName, ownerName).Scan(&roomID)
}

type ChatSearch struct {
	Owner               string
	Viewer              string
	ShowOnlyUnreadChats bool
}

func (n *ChatSearch) GetQuery() string {
	query := "WHERE"
	if n.Owner != "" {
		query += " owner_name = $1"
	} else {
		query += " account_name = $1"
	}

	if n.ShowOnlyUnreadChats {
		query += " AND new_message = true"
	}

	return query
}

func (n *ChatSearch) GetViewer() string {
	if n.Owner != "" {
		return n.Owner
	}
	return n.Viewer
}

func GetRoomsPageForwards(amount int, timeStamp time.Time, memberName string, info *ChatSearch) ([]ChatRoom, error) {
	result, err := postgresDB.Query(`SELECT room_id, member, account_name, new_message FROM chat_rooms_linked `+
		info.GetQuery()+` AND (created, account_name) <= ($2, $3) ORDER BY (created, account_name) DESC LIMIT $4`,
		info.GetViewer(), timeStamp, memberName, amount+1)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]ChatRoom, 0)
	chat := ChatRoom{}
	for result.Next() {
		err = result.Scan(&chat.Name, pq.Array(&chat.Member), &chat.User, &chat.NewMessage)
		if err != nil {
			return nil, err
		}
		arr = append(arr, chat)
	}
	return arr, nil
}

func GetRoomsPageBackwards(amount int, timeStamp time.Time, memberName string, info *ChatSearch) ([]ChatRoom, error) {
	result, err := postgresDB.Query(`SELECT room_id, member, account_name, new_message FROM (
SELECT room_id, member, account_name, new_message, created FROM chat_rooms_linked `+
		info.GetQuery()+` AND (created, account_name) >= ($2, $3) ORDER BY (created, account_name) LIMIT $4) 
    as room ORDER BY (created, account_name) DESC`, info.GetViewer(), timeStamp, memberName, amount+2)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]ChatRoom, 0)
	chat := ChatRoom{}
	for result.Next() {
		err = result.Scan(&chat.Name, pq.Array(&chat.Member), &chat.User, &chat.NewMessage)
		if err != nil {
			return nil, err
		}
		arr = append(arr, chat)
	}
	return arr, nil
}

func SetUnreadMessages(roomID string, viewer []string) {
	_, _ = postgresDB.Exec(`UPDATE chat_rooms_to_account SET new_message = true WHERE room_id = $1 AND (NOT (account_name = ANY($2)))`,
		roomID, pq.Array(viewer))
}

func SetReadMessage(roomID string, user string) {
	_, _ = postgresDB.Exec(`UPDATE chat_rooms_to_account SET new_message = false WHERE room_id = $1 AND account_name = $2`, roomID, user)
}
