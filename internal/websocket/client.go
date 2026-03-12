package websocket

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"chat-server/internal/service"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	userID uuid.UUID
	roomID uuid.UUID
	conn   *gorillaws.Conn
	send   chan []byte
	room   *Room
}

func NewClient(userID, roomID uuid.UUID, conn *gorillaws.Conn, room *Room) *Client {
	return &Client{
		userID: userID,
		roomID: roomID,
		conn:   conn,
		send:   make(chan []byte, 256),
		room:   room,
	}
}

func (c *Client) UserID() uuid.UUID { return c.userID }
func (c *Client) RoomID() uuid.UUID { return c.roomID }

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) ReadPump(ms *service.MessageService, ur *repository.UserRepository) {
	defer func() {
		ur.UpdateStatus(c.userID, models.UserStatusOffline)
		c.room.leave <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingEvent
		if err := json.Unmarshal(data, &incoming); err != nil {
			continue
		}

		switch incoming.Type {
		case EventMessage:
			msg, err := ms.CreateMessage(c.roomID, c.userID, incoming.Content, models.MessageTypeText, c.userID)
			if err != nil {
				continue
			}
			out, _ := json.Marshal(Event{Type: EventMessage, Payload: msg})
			c.room.broadcast <- out
		case EventHistory:
			msg, err := ms.ListMessages(c.roomID, incoming.Cursor, c.userID)
			if err != nil {
				continue
			}
			out, _ := json.Marshal(Event{Type: EventHistory, Payload: msg})
			c.send <- out
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(gorillaws.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(gorillaws.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(gorillaws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
