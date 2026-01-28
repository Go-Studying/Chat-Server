package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleMember MemberRole = "member"
)

type RoomMember struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RoomID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_room_user,unique" json:"room_id"`
	UserID   uuid.UUID  `gorm:"type:uuid;not_null;index:idx_room_user,unique" json:"user_id"`
	Role     MemberRole `gorm:"size:20;not null;default:member" json:"role"`
	JoinedAt time.Time  `gorm:"autoCreateTime" json:"joined_at"`

	// Relations
	Room ChatRoom `gorm:"foreignKey:RoomID" json:"-"`
	User User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (RoomMember) TableName() string {
	return "room_members"
}

func (rm *RoomMember) BeforeCreate(tx *gorm.DB) error {
	if rm.ID == uuid.Nil {
		rm.ID = uuid.New()
	}
	return nil
}
