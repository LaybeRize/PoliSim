package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	loc "PoliSim/localisation"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	id    string
	owner *database.Account
	hub   *Hub
	conn  *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 2500
)

func init() {
	go func() {
		for {
			select {
			case blocked, ok := <-database.BlockedAccountChannel:
				if !ok {
					return
				}
				HubMutex.Lock()
				for _, hub := range HubList {
					hub.Lock()
					for client := range hub.clients {
						if client.id == blocked {
							hub.unregister <- client
						}
					}
					hub.Unlock()
				}
				HubMutex.Unlock()
			case arr, ok := <-database.OwnerChangeOnAccountChannel:
				if !ok {
					return
				}
				HubMutex.Lock()
				for _, hub := range HubList {
					hub.Lock()
					for client := range hub.clients {
						if client.id == arr[1] && client.owner.Name != arr[0] {
							hub.unregister <- client
						}
					}
					hub.Unlock()
				}
				HubMutex.Unlock()
			}
		}
	}()
}

func serveWs(hub *Hub, owner *database.Account, user string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), id: user, owner: owner}

	go client.writePump()
	go client.readPump()
	client.hub.register <- client
}

const loadMessageAmount = 20

// Reads from Hub to the websocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	messages, err := database.LoadLastMessages(loadMessageAmount+1, time.Now().UTC().Add(time.Hour), c.hub.id, c.id)
	if err == nil {
		msgLen := len(messages)
		if len(messages) > loadMessageAmount {
			msgLen = loadMessageAmount
			c.send <- getButtonUpdate(c.hub.id, c.id, messages[loadMessageAmount-1].SendDate)
		}
		for i := range msgLen {
			if msgLen-i-1 < 0 {
				break
			}
			c.send <- getMessageTemplate(&messages[msgLen-i-1], c.owner, c.id)
		}

	}

	var w io.WriteCloser
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err = c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(msg)

			if err = w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err = c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Reads from WS connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()
		slog.Debug("received:", "user", c.id, "text", text)
		if err != nil {
			slog.Debug("error:", "error", err.Error())
			break
		}

		msg := &struct {
			Text string `json:"text"`
		}{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			slog.Debug("error:", "error", err.Error())
			continue
		}

		msg.Text = strings.TrimSpace(msg.Text)

		if len([]rune(msg.Text)) > 2000 || msg.Text == "" {
			continue
		}

		htmlMsg := &database.Message{
			SenderName: c.id,
			SendDate:   time.Now().UTC(),
			Text:       template.HTML(strings.ReplaceAll(template.HTMLEscaper(msg.Text), "\n", "<br>")),
		}
		err = database.InsertMessage(htmlMsg, c.hub.id)
		if err != nil {
			slog.Debug("error:", "error", err.Error())
			continue
		}
		c.hub.broadcast <- htmlMsg
		c.send <- []byte(loc.ChatRoomMessageTextarea)

	}
}

type Hub struct {
	sync.RWMutex

	id      string
	clients map[*Client]bool

	broadcast   chan *database.Message
	register    chan *Client
	unregister  chan *Client
	messageSend bool
}

func NewHub(id string) *Hub {
	return &Hub{
		id:          id,
		clients:     map[*Client]bool{},
		broadcast:   make(chan *database.Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		messageSend: false,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()
			database.SetReadMessage(h.id, client.id)
			slog.Debug("client registered", "id", client.id)
		case client := <-h.unregister:
			h.Lock()
			if h.messageSend {
				h.messageSend = false
			}
			if _, ok := h.clients[client]; ok {
				close(client.send)
				slog.Debug("client unregistered", "id", client.id)
				delete(h.clients, client)
			}
			if len(h.clients) == 0 {
				HubMutex.Lock()
				delete(HubList, h.id)
				HubMutex.Unlock()
				return
			}
			h.Unlock()
		case msg := <-h.broadcast:
			h.RLock()
			if !h.messageSend {
				h.messageSend = true
				go database.SetUnreadMessages(h.id, h.getCurrentViewers())
			}
			for client := range h.clients {
				select {
				case client.send <- getMessageTemplate(msg, client.owner, client.id):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.RUnlock()
		}
	}
}

func (h *Hub) getCurrentViewers() []string {
	res := make([]string, 0)
	for c := range h.clients {
		res = append(res, c.id)
	}
	return res
}

func getButtonUpdate(roomID string, receiverID string, timeStamp time.Time) []byte {
	var renderedMessage bytes.Buffer
	err := handler.MakeSpecialPagePartForWriter(&renderedMessage, &handler.ChatButtonObject{
		Room:      roomID,
		NextTime:  timeStamp,
		Recipient: receiverID,
	})
	if err != nil {
		log.Printf("error parsing Message: %v\n", err)
	}

	return renderedMessage.Bytes()
}

func getMessageTemplate(msg *database.Message, acc *database.Account, receiverID string) []byte {
	var renderedMessage bytes.Buffer
	err := handler.MakeSpecialPagePartForWriter(&renderedMessage, &handler.ChatMessageObject{
		Msg:       msg,
		Account:   acc,
		Recipient: receiverID,
	})
	if err != nil {
		log.Printf("error parsing Message: %v\n", err)
	}

	return renderedMessage.Bytes()
}
