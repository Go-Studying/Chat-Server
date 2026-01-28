package models

type ChatRoom struct {
	Base
	Name string `gorm:"size:100" json:"name"`

	// Relations
	Members  []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
	Messages []Message    `gorm:"foreignKey:RoomID" json:"-"`
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}
