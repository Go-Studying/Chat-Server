package websocket

import (
	"github.com/google/uuid"
)

type Room struct {
	id        uuid.UUID
	clients   map[uuid.UUID]*Client
	broadcast chan []byte
	join      chan *Client
	leave     chan *Client
}

func NewRoom(id uuid.UUID) *Room {
	return &Room{
		id:        id,
		clients:   make(map[uuid.UUID]*Client),
		broadcast: make(chan []byte, 256),
		join:      make(chan *Client),
		leave:     make(chan *Client),
	}
}

func (r *Room) Join(client *Client) {
	r.join <- client
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client.userID] = client

		case client := <-r.leave:
			if _, ok := r.clients[client.userID]; ok {
				delete(r.clients, client.userID)
				close(client.send)
			}

		case message := <-r.broadcast:
			for userID, client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, userID)
				}
			}
		}
	}
}
