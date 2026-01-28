package models

import "github.com/google/uuid"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	Base
	RoomID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"room_id"`
	SenderID uuid.UUID   `gorm:"type:uuid;not null;index" json:"sender_id"`
	Content  string      `gorm:"type:text;not null" json:"content"`
	Type     MessageType `gorm:"size:20;not null;default:text" json:"type"`

	// Relations
	Room   ChatRoom `gorm:"foreignKey:RoomID" json:"-"`
	Sender User     `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (Message) Tablename() string {
	return "messages"
}
