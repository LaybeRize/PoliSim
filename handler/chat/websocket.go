package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"log/slog"
	"net/http"
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
	maxMessageSize = 512
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
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// Reads from Hub to the websocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	messages, err := database.LoadLastMessages(20, time.Now().UTC().Add(time.Hour), c.hub.id, c.id)
	if err == nil {
		for _, msg := range messages {
			c.send <- getMessageTemplate(&msg, c.owner, c.id)
		}
		// Todo send an update for the Load next messages button
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

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(msg)
			}

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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Debug("error:", "error", err.Error())
			}
			break
		}

		msg := &database.Message{}

		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			slog.Debug("error:", "error", err.Error())
		}

		msg.SenderName = c.id
		msg.SendDate = time.Now().UTC()
		c.hub.broadcast <- msg
	}
}

type Hub struct {
	sync.RWMutex

	id      string
	clients map[*Client]bool

	broadcast  chan *database.Message
	register   chan *Client
	unregister chan *Client
}

func NewHub(id string) *Hub {
	return &Hub{
		id:         id,
		clients:    map[*Client]bool{},
		broadcast:  make(chan *database.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()

			slog.Debug("client registered", "id", client.id)
		case client := <-h.unregister:
			h.Lock()
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
			err := database.InsertMessage(msg, h.id)
			if err != nil {
				continue
			}
			h.RLock()
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
