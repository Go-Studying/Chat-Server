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

func (s *ChatRoomService) JoinRoom(id, userID uuid.UUID) error {
	room, err := s.r.FindByID(id)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}
	for _, m := range room.Members {
		if m.UserID == userID {
			return ErrMemberExists
		}
	}

	member := models.RoomMember{
		RoomID: room.ID,
		UserID: userID,
		Role:   models.MemberRoleMember,
	}
	return s.r.AddMember(room.ID, &member)
}

func (s *ChatRoomService) GetRoom(id, userID uuid.UUID) (*models.ChatRoom, error) {
	room, err := s.r.FindByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}
	isMember := false
	for _, m := range room.Members {
		if m.UserID == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, ErrNotMember
	}
	return room, nil
}

func (s *ChatRoomService) GetMyRooms(userID uuid.UUID) ([]models.ChatRoom, error) {
	return s.r.FindByUserID(userID)
}

func (s *ChatRoomService) DeleteRoom(roomID, userID uuid.UUID) error {
	room, err := s.r.FindByID(roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}
	if room.CreatedBy != userID {
		return ErrNotOwner
	}
	return s.r.Delete(roomID)
}
