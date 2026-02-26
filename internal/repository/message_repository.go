package repository

import (
	"chat-server/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) FindAll(roomID uuid.UUID, cursor *time.Time) ([]models.Message, error) {
	var messages []models.Message

	query := r.db.Where("room_id = ?", roomID)
	if cursor != nil {
		query = query.Where("created_at < ?", cursor)
	}

	err := r.db.
		Order("created_at DESC").
		Limit(10).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
