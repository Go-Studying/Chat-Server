package main

import (
	"chat-server/internal/config"
	"chat-server/internal/routes"
	"log"
)

func main() {
	// 설정 로드
	cfg := config.Load()

	// 라우터 설정
	router := routes.SetupRouter(cfg)

	// 서버 시작
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.Environment)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
