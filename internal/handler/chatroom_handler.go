package handler

import (
	"chat-server/internal/middleware"
	"chat-server/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

type ChatRoomHandler struct {
	chatRoomService *service.ChatRoomService
	sanitizer       *bluemonday.Policy
}

func NewChatRoomHandler(s *service.ChatRoomService) *ChatRoomHandler {
	return &ChatRoomHandler{
		chatRoomService: s,
		sanitizer:       bluemonday.StrictPolicy(),
	}
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

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	sanitizedName := h.sanitizer.Sanitize(req.Name)
	if sanitizedName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room name"})
		return
	}

	room, err := h.chatRoomService.CreateRoom(sanitizedName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *ChatRoomHandler) GetMyRooms(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

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

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	room, err := h.chatRoomService.GetRoom(roomID, userID)
	if err != nil {
		if errors.Is(err, service.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		if errors.Is(err, service.ErrNotMember) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
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

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	if err := h.chatRoomService.DeleteRoom(roomID, userID); err != nil {
		if errors.Is(err, service.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete room"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}
