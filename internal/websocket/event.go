package websocket

import "time"

type EventType string

const (
	EventMessage EventType = "message"
	EventHistory EventType = "history"
	EventJoin    EventType = "join"
	EventLeave   EventType = "leave"
	EventOnline  EventType = "user_online"
	EventOffline EventType = "user_offline"
)

type Event struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type IncomingEvent struct {
	Type    EventType  `json:"type"`
	Content string     `json:"content,omitempty"`
	Cursor  *time.Time `json:"cursor,omitempty"`
}
