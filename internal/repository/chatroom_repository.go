package repository

import (
	"chat-server/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRoomRepository struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) *ChatRoomRepository {
	return &ChatRoomRepository{db: db}
}

func (r *ChatRoomRepository) Create(room *models.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *ChatRoomRepository) FindByID(id uuid.UUID) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.Preload("Members.User").Take(&room, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *ChatRoomRepository) FindByUserID(userID uuid.UUID) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := r.db.
		Joins("JOIN room_members ON room_members.room_id = chat_rooms.id").
		Where("room_members.user_id = ? AND room_members.deleted_at IS NULL", userID).
		Preload("Members.User").
		Find(&rooms).Error
	return rooms, err
}

func (r *ChatRoomRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ChatRoom{}, "id = ?", id).Error
}

func (r *ChatRoomRepository) AddMember(roomID uuid.UUID, member *models.RoomMember) error {
	return r.db.Model(&models.ChatRoom{Base: models.Base{ID: roomID}}).
		Association("Members").
		Append(member)
}
