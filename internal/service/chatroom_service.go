package service

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"

	"github.com/google/uuid"
)

type ChatRoomService struct {
	r *repository.ChatRoomRepository
}

func NewChatRoomService(r *repository.ChatRoomRepository) *ChatRoomService {
	return &ChatRoomService{r: r}
}

func (s *ChatRoomService) CreateRoom(name string, ownerID uuid.UUID) (*models.ChatRoom, error) {
	room := &models.ChatRoom{
		Name:      name,
		CreatedBy: ownerID,
		Members: []models.RoomMember{
			{UserID: ownerID, Role: models.MemberRoleOwner},
		},
	}
	if err := s.r.Create(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *ChatRoomService) GetRoom(id uuid.UUID) (*models.ChatRoom, error) {
	room, err := s.r.FindByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *ChatRoomService) GetMyRooms(userID uuid.UUID) ([]models.ChatRoom, error) {
	return s.r.FindByUserID(userID)
}

func (s *ChatRoomService) DeleteRoom(roomID uuid.UUID) error {
	return s.r.Delete(roomID)
}
