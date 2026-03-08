package chat

import (
	"chat-server/internal/middleware"
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"chat-server/internal/service"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 중앙 관리자 역할
type Hub struct {
	//1. Hub.Run() goroutine + rooms 채널 관리
	//rooms chan *Room

	//2. rooms를 map으로 관리 (동시성 보장x) + mutex 이용
	rooms          map[uuid.UUID]*Room
	mutex          sync.RWMutex
	userRepository *repository.UserRepository
	messageService *service.MessageService
}

func NewHub(ur *repository.UserRepository, ms *service.MessageService) *Hub {
	return &Hub{
		//rooms: make(chan *Room),
		rooms:          make(map[uuid.UUID]*Room),
		userRepository: ur,
		messageService: ms,
	}
}

func (h *Hub) ServeWs(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	room := h.GetOrCreateRoom(roomID)

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	user, err := h.userRepository.FindUserByID(userID)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //get a websocket connection
	client := NewClient(userID, user.Username, room, conn)
	client.room.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) GetOrCreateRoom(roomID uuid.UUID) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if room, ok := h.rooms[roomID]; ok {
		return room
	}

	room := NewRoom(h, roomID)
	h.rooms[roomID] = room
	go room.Run()

	return room
}

func (h *Hub) SaveMessage(roomID uuid.UUID, senderID uuid.UUID, content string) error {
	_, err := h.messageService.CreateMessage(roomID, senderID, content, models.MessageTypeText, senderID)
	return err
}
