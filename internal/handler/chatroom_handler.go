package handler

import (
	"chat-server/internal/middleware"
	"chat-server/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatRoomHandler struct {
	chatRoomService *service.ChatRoomService
}

func NewChatRoomHandler(s *service.ChatRoomService) *ChatRoomHandler {
	return &ChatRoomHandler{chatRoomService: s}
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

func (h *ChatRoomHandler) Create(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userIDStr, _ := middleware.GetCurrentUser(c)
	userID, _ := uuid.Parse(userIDStr)

	room, err := h.chatRoomService.CreateRoom(req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *ChatRoomHandler) GetMyRooms(c *gin.Context) {
	userIDStr, _ := middleware.GetCurrentUser(c)
	userID, _ := uuid.Parse(userIDStr)

	rooms, err := h.chatRoomService.GetMyRooms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *ChatRoomHandler) GetRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.chatRoomService.GetRoom(roomID)
	if err != nil {
		if errors.Is(err, service.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get room"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *ChatRoomHandler) Delete(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	if err := h.chatRoomService.DeleteRoom(roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}
