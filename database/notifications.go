package database

import (
	"log"
	"strconv"
	"sync/atomic"
)

var createdAccountForNotifications = make(chan string)
var blockedAccountChannelForNotifications = make(chan string)
var ownerChangeOnAccountChannelForNotifications = make(chan []string)

var markLetterAsRead = make(chan []string)
var markChatAsRead = make(chan []string)
var newLetter = make(chan newEntry)
var newChatMessage = make(chan newEntry)

type newEntry struct {
	ID       string
	Accounts []string
}

type Notification struct {
	HasUnreadLetters bool
	HasUnreadChats   bool
	UnreadLetters    string
	UnreadChats      string
}

func (n *Notification) HasNotifications() bool {
	if n == nil {
		return false
	}
	return n.HasUnreadLetters || n.HasUnreadChats
}

func GetUserNotification(acc *Account) *Notification {
	if acc == nil {
		return nil
	}
	val := notificationMap[acc.Name]
	result := &Notification{
		HasUnreadLetters: false,
		HasUnreadChats:   false,
	}
	if val == nil {
		return result
	}
	if letter := int(val.UnreadLetter.Load()); letter < 10 && letter > 0 {
		result.HasUnreadLetters = true
		result.UnreadLetters = strconv.Itoa(letter)
	} else if letter > 0 {
		result.HasUnreadLetters = true
		result.UnreadLetters = "9+"
	}

	if chat := int(val.UnreadChats.Load()); chat < 10 && chat > 0 {
		result.HasUnreadChats = true
		result.UnreadChats = strconv.Itoa(chat)
	} else if chat > 0 {
		result.HasUnreadChats = true
		result.UnreadChats = "9+"
	}
	return result
}

type notification struct {
	Parent          *notification
	UnreadLetter    atomic.Int64
	UnreadLetterMap map[string]struct{}
	UnreadChats     atomic.Int64
	UnreadChatsMap  map[string]struct{}
}

var notificationMap = make(map[string]*notification)

func init() {
	go func() {
		for {
			select {
			// New Account was created
			case val, ok := <-createdAccountForNotifications:
				if !ok {
					return
				}
				notificationMap[val] = &notification{
					Parent:          nil,
					UnreadLetter:    atomic.Int64{},
					UnreadLetterMap: make(map[string]struct{}),
					UnreadChats:     atomic.Int64{},
					UnreadChatsMap:  make(map[string]struct{}),
				}
			case val, ok := <-blockedAccountChannelForNotifications:
				if !ok {
					return
				}
				markAccountAsBlocked(val)
			case val, ok := <-ownerChangeOnAccountChannelForNotifications:
				if !ok {
					return
				}
				accountChangedOwner(val[0], val[1])
			case val, ok := <-markLetterAsRead:
				if !ok {
					return
				}
				markLetterForAccountAsRead(val[0], val[1])
			case val, ok := <-markChatAsRead:
				if !ok {
					return
				}
				markChatMessagesAsRead(val[0], val[1])
			case val, ok := <-newLetter:
				if !ok {
					return
				}
				usersHaveNewUnreadLetter(val)
			case val, ok := <-newChatMessage:
				if !ok {
					return
				}
				usersHaveNewUnreadChatMessage(val)
			}
		}
	}()
}

func markAccountAsBlocked(account string) {
	if notificationMap[account].Parent == nil {
		return
	}
	notif := notificationMap[account]
	parent := notif.Parent
	if parent != nil {
		parent.UnreadLetter.Add(-notif.UnreadLetter.Load())
		parent.UnreadChats.Add(-notif.UnreadChats.Load())
	}
	notif.Parent = nil
}

func accountChangedOwner(owner string, account string) {
	if notificationMap[account].Parent != nil {
		notif := notificationMap[account]
		parent := notif.Parent

		parent.UnreadLetter.Add(-notif.UnreadLetter.Load())
		parent.UnreadChats.Add(-notif.UnreadChats.Load())
	}

	if owner == account {
		notificationMap[account].Parent = nil
	} else {
		notificationMap[account].Parent = notificationMap[owner]
		notif := notificationMap[account]
		parent := notif.Parent

		parent.UnreadLetter.Add(notif.UnreadLetter.Load())
		parent.UnreadChats.Add(notif.UnreadChats.Load())
	}
}

func markLetterForAccountAsRead(account string, letterID string) {
	notif := notificationMap[account]
	if notif == nil {
		return
	}
	if _, exists := notif.UnreadLetterMap[letterID]; exists {
		notif.UnreadLetter.Add(-1)
		if notif.Parent != nil {
			notif.Parent.UnreadLetter.Add(-1)
		}
		delete(notif.UnreadLetterMap, letterID)
	}
}

func markChatMessagesAsRead(account string, roomID string) {
	notif := notificationMap[account]
	if notif == nil {
		return
	}
	if _, exists := notif.UnreadChatsMap[roomID]; exists {
		notif.UnreadChats.Add(-1)
		if notif.Parent != nil {
			notif.Parent.UnreadChats.Add(-1)
		}
		delete(notif.UnreadChatsMap, roomID)
	}
}

func usersHaveNewUnreadLetter(val newEntry) {
	for _, name := range val.Accounts {
		notif := notificationMap[name]
		if notif == nil {
			continue
		}
		parent := notif.Parent

		notif.UnreadLetter.Add(1)
		notif.UnreadLetterMap[val.ID] = struct{}{}
		if parent != nil {
			parent.UnreadLetter.Add(1)
		}
	}
}

func usersHaveNewUnreadChatMessage(val newEntry) {
	for _, name := range val.Accounts {
		notif := notificationMap[name]
		if notif == nil {
			continue
		}
		if _, exists := notif.UnreadChatsMap[val.ID]; exists {
			continue
		}
		parent := notif.Parent

		notif.UnreadChats.Add(1)
		notif.UnreadChatsMap[val.ID] = struct{}{}
		if parent != nil {
			parent.UnreadChats.Add(1)
		}
	}
}

func loadNotificationsFromDB() {
	result, err := postgresDB.Query("SELECT account_name, owner_name FROM ownership;")
	if err != nil {
		log.Fatalf("Could not read postgres ownership table: %v", err)
	}

	for result.Next() {
		var account string
		var owner string
		err = result.Scan(&account, &owner)
		if err != nil {
			log.Fatalf("Could not scan entry from postgres ownership table: %v", err)
		}
		addNewOwnerAccountPair(owner, account)
	}

	result, err = postgresDB.Query("SELECT account_name, letter_id FROM letter_to_account WHERE has_read = false;")
	if err != nil {
		log.Fatalf("Could not read postgres letter_to_account table: %v", err)
	}
	for result.Next() {
		var account string
		var letterID string
		err = result.Scan(&account, &letterID)
		if err != nil {
			log.Fatalf("Could not scan entry from postgres letter_to_account table: %v", err)
		}
		usersHaveNewUnreadLetter(newEntry{ID: letterID, Accounts: []string{account}})
	}

	result, err = postgresDB.Query("SELECT account_name, room_id FROM chat_rooms_to_account WHERE new_message = true;")
	if err != nil {
		log.Fatalf("Could not read postgres chat_rooms_to_account table: %v", err)
	}
	for result.Next() {
		var account string
		var roomID string
		err = result.Scan(&account, &roomID)
		if err != nil {
			log.Fatalf("Could not scan entry from postgres chat_rooms_to_account table: %v", err)
		}
		usersHaveNewUnreadChatMessage(newEntry{ID: roomID, Accounts: []string{account}})
	}
}

func addNewOwnerAccountPair(owner string, account string) {
	if _, exists := notificationMap[account]; !exists {
		notificationMap[account] = &notification{
			Parent:          nil,
			UnreadLetter:    atomic.Int64{},
			UnreadLetterMap: make(map[string]struct{}),
			UnreadChats:     atomic.Int64{},
			UnreadChatsMap:  make(map[string]struct{}),
		}
	}
	if _, exists := notificationMap[owner]; !exists {
		notificationMap[owner] = &notification{
			Parent:          nil,
			UnreadLetter:    atomic.Int64{},
			UnreadLetterMap: make(map[string]struct{}),
			UnreadChats:     atomic.Int64{},
			UnreadChatsMap:  make(map[string]struct{}),
		}
	}
	if owner != account {
		notificationMap[account].Parent = notificationMap[owner]
	}
}
