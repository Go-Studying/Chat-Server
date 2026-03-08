package chat

import "github.com/google/uuid"

type Room struct {
	ID           uuid.UUID
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	messageSaver MessageSaver
}

func NewRoom(messageSaver MessageSaver, id uuid.UUID) *Room {
	return &Room{
		ID:           id,
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		clients:      make(map[*Client]bool),
		messageSaver: messageSaver,
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

type MessageSaver interface {
	SaveMessage(roomID, senderID uuid.UUID, contents string) error
}
