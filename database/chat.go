package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Message struct {
	SenderName string    `json:"-"`
	SendDate   time.Time `json:"-"`
	Text       string    `json:"text"`
}

type ChatRoom struct {
	Name   string
	Member []string
}

func (m *Message) GetTimeSend(a *Account) string {
	return m.SendDate.In(a.TimeZone).Format(loc.TimeFormatString)
}

func LoadLastMessages(amount int, timeStamp time.Time, roomID string, accountName string) ([]Message, error) {
	err := postgresDB.QueryRow(`SELECT account_name from chat_rooms_to_account WHERE account_name = $1 AND room_id = $2`, accountName, roomID).Scan(&accountName)
	if err != nil {
		return nil, err
	}
	result, err := postgresDB.Query(`SELECT sender, send_time, message FROM 
(SELECT sender, send_time, message FROM chat_messages WHERE send_time < $1 AND room_id = $2 ORDER BY send_time DESC LIMIT $3) as msg 
ORDER BY msg.send_time`, timeStamp, roomID, amount)
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
	err = tx.QueryRow(`SELECT room_id FROM chat_rooms WHERE member = ARRAY(SELECT name FROM account WHERE name = ANY($1) ORDER BY name)`, pq.Array(member)).Scan(name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err == nil {
		return DoubleChatRoomEntry
	}

	_, err = tx.Exec(`INSERT INTO chat_rooms (room_id, member) VALUES ($1, ARRAY(SELECT name FROM account WHERE name = ANY($2) ORDER BY name))`, roomID, pq.Array(member))
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO chat_rooms_to_account (room_id, account_name) 
SELECT $1 AS room_id, name FROM account WHERE name = ANY($2)`, roomID, pq.Array(member))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func QueryForRoomIdAndUser(roomID string, accountName string, ownerName string) error {
	return postgresDB.QueryRow(`SELECT room_id FROM chat_rooms_to_account cta 
    INNER JOIN ownership own ON cta.account_name = own.account_name WHERE room_id = $1 AND own.account_name = $2 AND owner_name = $3;`,
		roomID, accountName, ownerName).Scan(&roomID)
}

func GetAllRoomsForUser(viewer []string) ([]ChatRoom, error) {
	result, err := postgresDB.Query(`SELECT DISTINCT ON (chat_rooms.room_id) chat_rooms.room_id, member FROM chat_rooms 
    INNER JOIN public.chat_rooms_to_account cta on chat_rooms.room_id = cta.room_id 
WHERE account_name = ANY($1) ORDER BY chat_rooms.room_id`, pq.Array(viewer))
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]ChatRoom, 0)
	chat := ChatRoom{}
	for result.Next() {
		err = result.Scan(&chat.Name, pq.Array(&chat.Member))
		if err != nil {
			return nil, err
		}
		arr = append(arr, chat)
	}
	return arr, nil
}
