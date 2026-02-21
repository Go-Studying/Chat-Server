package routes

import (
	"chat-server/internal/config"
	"chat-server/internal/handler"
	"chat-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler, chatRoomHandler *handler.ChatRoomHandler) *gin.Engine {
	// 환경에 따른 Gin 모드 설정
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS 미들웨어(필요하면 작성)

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Chat server is running",
		})
	})

	// 인증
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", authHandler.SignUp)
		auth.POST("/login", authHandler.Login)
	}

	// 사용자 인증 필요한 api
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		rooms := api.Group("/rooms")
		{
			rooms.POST("", chatRoomHandler.Create)
			rooms.GET("", chatRoomHandler.GetMyRooms)
			rooms.GET("/:id", chatRoomHandler.GetRoom)
			rooms.DELETE("/:id", chatRoomHandler.Delete)
		}
	}

	return router
}
