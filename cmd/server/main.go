package main

import (
	"chat-server/internal/config"
	"chat-server/internal/handler"
	"chat-server/internal/models/database"
	"chat-server/internal/repository"
	"chat-server/internal/routes"
	"chat-server/internal/service"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 설정 로드
	cfg := config.Load()

	// DB 연결
	db, err := database.New(cfg)
	if err != nil {
		slog.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 의존성 주입
	userRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, cfg)

	chatRoomRepo := repository.NewChatRoomRepository(db.DB)
	chatRoomService := service.NewChatRoomService(chatRoomRepo)
	chatRoomHandler := handler.NewChatRoomHandler(chatRoomService)

	// 마이그레이션
	if err := db.AutoMigrate(); err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migration completed")

	// 라우터 설정
	router := routes.SetupRouter(cfg, authHandler, chatRoomHandler)

	// 서버 시작
	addr := ":" + cfg.Port
	slog.Info("Server starting", "address", addr)
	slog.Info("Environment", "env", cfg.Environment)

	if err := router.Run(addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
