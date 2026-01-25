package main

import (
	"chat-server/internal/config"
	"chat-server/internal/routes"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 설정 로드
	cfg := config.Load()

	// 라우터 설정
	router := routes.SetupRouter(cfg)

	// 서버 시작
	addr := ":" + cfg.Port
	slog.Info("Server starting", "address", addr)
	slog.Info("Environment", "env", cfg.Environment)

	if err := router.Run(addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
