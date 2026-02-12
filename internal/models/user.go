package models

type UserStatus string

const (
	UserStatusOnline  UserStatus = "online"
	UserStatusOffline UserStatus = "offline"
)

type User struct {
	Base
	Email    string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Username string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password string     `gorm:"size:255;not null" json:"-"` // json응답에서 제외
	Status   UserStatus `gorm:"size:20;default:online" json:"status"`

	// Relation
	Messages    []Message    `gorm:"foreignKey:SenderID" json:"-"`
	RoomMembers []RoomMember `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}
