package service

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"time"

	"github.com/google/uuid"
)

type MessageService struct {
	mr *repository.MessageRepository
	cs *ChatRoomService
}

func NewMessageService(mr *repository.MessageRepository, cs *ChatRoomService) *MessageService {
	return &MessageService{mr: mr, cs: cs}
}

func (s *MessageService) CreateMessage(roomID uuid.UUID, senderID uuid.UUID, content string, messageType models.MessageType, userID uuid.UUID) (*models.Message, error) {
	_, err := s.cs.GetRoom(roomID, userID)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
		Type:     messageType,
	}

	if err := s.mr.Create(message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) ListMessages(roomID uuid.UUID, cursor *time.Time, userID uuid.UUID) ([]models.Message, error) {
	room, err := s.cs.GetRoom(roomID, userID)
	if err != nil {
		return nil, err
	}

	return s.mr.FindAll(room.ID, cursor)
}
