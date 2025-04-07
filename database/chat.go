package database

import (
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"html/template"
	"net/url"
	"strconv"
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
	ID         string
	Member     []string
	User       string
	NewMessage bool
	Created    time.Time
}

func (c *ChatRoom) GetLink() template.URL {
	return template.URL("/chat/" + c.ID + "/" + url.PathEscape(c.User))
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

func CreateChatRoom(roomName string, member []string) error {
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

	roomID := helper.GetUniqueID(roomName)
	_, err = tx.Exec(`INSERT INTO chat_rooms (name, member, created, room_id) VALUES ($1, 
	ARRAY(SELECT name FROM account WHERE name = ANY($2) AND name <> $3 AND blocked = false ORDER BY name), $4, $5)`,
		roomName, pq.Array(member), loc.AdministrationAccountName, time.Now().UTC(), roomID)
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
	Name                string
	Owner               string
	Viewer              string
	ShowOnlyUnreadChats bool
	vals                []any
}

func (n *ChatSearch) GetQuery() string {
	n.vals = make([]any, 1, 3)
	n.vals[0] = n.Owner
	query := " WHERE owner_name = $4"
	pos := 5

	if n.Viewer != "" {
		n.vals = append(n.vals, n.Viewer)
		query += " AND account_name = $" + strconv.Itoa(pos) + " "
		pos += 1
	}

	if n.ShowOnlyUnreadChats {
		query += " AND new_message = true "
	}

	if n.Name != "" {
		n.vals = append(n.vals, n.Name)
		query += " AND name LIKE '%' || $" + strconv.Itoa(pos) + " || '%' "
		pos += 1
	}

	return query
}

func (n *ChatSearch) GetViewer(in []any) []any {
	return append(in, n.vals...)
}

func GetRoomsPageForwards(amount int, timeStamp time.Time, memberName string, info *ChatSearch) ([]ChatRoom, error) {
	result, err := postgresDB.Query(`SELECT room_id, member, account_name, new_message, created, name FROM chat_rooms_linked `+
		info.GetQuery()+` AND (created, account_name) <= ($1, $2) ORDER BY (created, account_name) DESC LIMIT $3`,
		info.GetViewer([]any{timeStamp, memberName, amount + 1})...)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]ChatRoom, 0)
	chat := ChatRoom{}
	for result.Next() {
		err = result.Scan(&chat.ID, pq.Array(&chat.Member), &chat.User, &chat.NewMessage, &chat.Created, &chat.Name)
		if err != nil {
			return nil, err
		}
		arr = append(arr, chat)
	}
	return arr, nil
}

func GetRoomsPageBackwards(amount int, timeStamp time.Time, memberName string, info *ChatSearch) ([]ChatRoom, error) {
	result, err := postgresDB.Query(`SELECT room_id, member, account_name, new_message, created, name FROM (
SELECT room_id, member, account_name, new_message, created, name FROM chat_rooms_linked `+
		info.GetQuery()+` AND (created, account_name) >= ($1, $2) ORDER BY (created, account_name) LIMIT $3) 
    as room ORDER BY (created, account_name) DESC`, info.GetViewer([]any{timeStamp, memberName, amount + 2}))
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]ChatRoom, 0)
	chat := ChatRoom{}
	for result.Next() {
		err = result.Scan(&chat.ID, pq.Array(&chat.Member), &chat.User, &chat.NewMessage, &chat.Created, &chat.Name)
		if err != nil {
			return nil, err
		}
		arr = append(arr, chat)
	}
	return arr, nil
}

func SetUnreadMessages(roomID string, viewer []string) {
	result, err := postgresDB.Query(`UPDATE chat_rooms_to_account SET new_message = true 
                             WHERE room_id = $1 AND (NOT (account_name = ANY($2))) RETURNING account_name;`,
		roomID, pq.Array(viewer))
	if err != nil {
		return
	}
	defer closeRows(result)
	val := newEntry{ID: roomID, Accounts: make([]string, 0)}
	for result.Next() {
		var name string
		err = result.Scan(&name)
		if err != nil {
			continue
		}
		val.Accounts = append(val.Accounts, name)
	}
	newChatMessage <- val
}

func SetReadMessage(roomID string, user string) {
	_, err := postgresDB.Exec(`UPDATE chat_rooms_to_account SET new_message = false WHERE room_id = $1 AND account_name = $2`, roomID, user)
	if err == nil {
		markChatAsRead <- []string{user, roomID}
	}
}
