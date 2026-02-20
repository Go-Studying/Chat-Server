package models

import "github.com/google/uuid"

type ChatRoom struct {
	Base
	Name      string    `gorm:"size:100" json:"name"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`

	// Relations
	Creator  User         `gorm:"foreignKey:CreatedBy" json:"-"`
	Members  []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
	Messages []Message    `gorm:"foreignKey:RoomID" json:"-"`
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}
