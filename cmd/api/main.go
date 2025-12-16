package main

import (
	"log"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/logger"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/server"
)

func main() {
	cfg := config.Load()
	log.Println(*cfg)

	logger.Setup(*cfg)

	handlers := handlers.New(nil)

	srv := server.New(*cfg, handlers)

	if err := srv.Start(); err != nil {
		slog.Error("Failed to start server",
			"error", err)
	}
}
