package routes

import (
	"chat-server/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
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

	// API 그룹
	/*api := router.Group("/api/")
	{
		// 나중에 루트 추가(API 구현시)
	}*/

	return router
}
