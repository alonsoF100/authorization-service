package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/router"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
	router *chi.Mux
	cfg    config.Config
}

func New(cfg config.Config, handlers *handlers.Handler) *Server {
	rtr := router.New(handlers).Setup()

	// Настраиваем HTTP сервер
	srv := &http.Server{
		Addr:         cfg.Server.PortStr(),
		Handler:      rtr,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.Server.IdleTimeout * time.Second,
	}

	return &Server{
		server: srv,
		router: rtr,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	slog.Info("Starting HTTP server",
		"address", s.cfg.Server.PortStr(),
		"read_timeout", s.cfg.Server.ReadTimeout,
		"write_timeout", s.cfg.Server.WriteTimeout,
		"idle_timeout", s.cfg.Server.IdleTimeout,
	)

	return s.server.ListenAndServe()
}
