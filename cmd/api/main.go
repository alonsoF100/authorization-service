package main

import (
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/logger"
	"github.com/alonsoF100/authorization-service/internal/repository/postgres"
	"github.com/alonsoF100/authorization-service/internal/service"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/server"
	_ "github.com/alonsoF100/authorization-service/migrations/postgres"
)

func main() {
	cfg := config.Load()

	logS := logger.Setup(cfg)

	pool, err := postgres.NewPool(cfg)
	if err != nil {
		slog.Error("Failed to create pool", "error", err)
	}
	defer pool.Close()
	slog.Info("Pool created successfully")

	dataBase := postgres.New(pool)

	authService := service.NewAuthService(
		dataBase,
		cfg,
	)
	userService := service.NewUserService(dataBase)

	handlers := handlers.New(
		authService,
		userService,
	)

	if err := server.New(cfg, handlers, logS).Start(); err != nil {
		slog.Error("Failed to start server",
			"error", err)
	}
}
