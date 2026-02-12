package handler

import (
	"chat-server/internal/config"
	"chat-server/internal/middleware"
	"chat-server/internal/service"
	"chat-server/internal/tools/security"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserID   uuid.UUID `json:"userId"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, err := h.authService.SignUp(req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign up"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userId": userID})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := security.NewJWT(userID, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetCookie(middleware.AuthCookieName,
		token,
		60*60*24,
		"/",
		"localhost",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"userId": userID})
}

func (h *AuthHandler) Test(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	c.JSON(http.StatusOK, gin.H{"hello": userID})
}
