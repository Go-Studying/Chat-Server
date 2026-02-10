package routes

import (
	"chat-server/internal/config"
	"chat-server/internal/handler"
	"chat-server/internal/middleware"
	"chat-server/internal/models/database"
	"chat-server/internal/repository"
	"chat-server/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	// 환경에 따른 Gin 모드 설정
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, _ := database.New(config.Load())
	userRepository := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

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
	api := router.Group("/api")
	{
		api.POST("/signup", authHandler.SignUp)
		api.POST("/login", authHandler.Login)
	}

	// 사용자 인증 필요한 api
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/test", authHandler.Test)
	}

	return router
}
