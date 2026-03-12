package handler

import (
	"chat-server/internal/config"
	"chat-server/internal/middleware"
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"chat-server/internal/service"
	"chat-server/internal/websocket"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader        gorillaws.Upgrader
	manager         *websocket.Manager
	messageService  *service.MessageService
	chatRoomService *service.ChatRoomService
	userRepo        *repository.UserRepository
}

func NewWebSocketHandler(cfg *config.Config, manager *websocket.Manager, ms *service.MessageService, cs *service.ChatRoomService, ur *repository.UserRepository) *WebSocketHandler {
	allowedOrigins := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedOrigins[o] = struct{}{}
	}

	return &WebSocketHandler{
		upgrader: gorillaws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return false
				}
				_, ok := allowedOrigins[origin]
				return ok
			},
		},
		manager:         manager,
		messageService:  ms,
		chatRoomService: cs,
		userRepo:        ur,
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

	// WebSocket 업그레이드 전에 멤버 확인
	if _, err := h.chatRoomService.GetRoom(roomID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	room := h.manager.GetOrCreateRoom(roomID)
	client := websocket.NewClient(userID, roomID, conn, room)
	room.Join(client)

	h.userRepo.UpdateStatus(userID, models.UserStatusOnline)

	go h.sendHistory(client)
	go client.WritePump()
	go client.ReadPump(h.messageService, h.userRepo)
}

func (h *WebSocketHandler) sendHistory(client *websocket.Client) {
	msgs, err := h.messageService.ListMessages(client.RoomID(), nil, client.UserID())
	if err != nil {
		return
	}
	out, err := json.Marshal(websocket.Event{Type: websocket.EventHistory, Payload: msgs})
	if err != nil {
		slog.Error("failed to marshal history event", "error", err)
		return
	}
	client.Send(out)
}
