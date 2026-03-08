package handler

import (
	"chat-server/internal/middleware"
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"chat-server/internal/service"
	"chat-server/internal/websocket"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	manager        *websocket.Manager
	messageService *service.MessageService
	userRepo       *repository.UserRepository
}

func NewWebSocketHandler(manager *websocket.Manager, ms *service.MessageService, ur *repository.UserRepository) *WebSocketHandler {
	return &WebSocketHandler{
		manager:        manager,
		messageService: ms,
		userRepo:       ur,
	}
}

func (h *WebSocketHandler) Connect(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	room := h.manager.GetOrCreateRoom(roomID)
	client := websocket.NewClient(userID, roomID, conn, room)
	room.Join(client)

	h.userRepo.UpdateStatus(userID, models.UserStatusOnline)

	go h.sendHistory(client)
	go client.WritePump()
	client.ReadPump(h.messageService, h.userRepo)
}

func (h *WebSocketHandler) sendHistory(client *websocket.Client) {
	msgs, err := h.messageService.ListMessages(client.RoomID(), nil, client.UserID())
	if err != nil {
		return
	}
	out, _ := json.Marshal(websocket.Event{Type: websocket.EventHistory, Payload: msgs})
	client.Send(out)
}
