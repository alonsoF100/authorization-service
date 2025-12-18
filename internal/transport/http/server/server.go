package server

import (
	"log/slog"
	"net/http"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/router"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Server *http.Server
	Router *chi.Mux
	Cfg    *config.Config
}

func New(cfg *config.Config, handlers *handlers.Handler) *Server {
	rtr := router.New(handlers).Setup()

	srv := &http.Server{
		Addr:         cfg.Server.PortStr(),
		Handler:      rtr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		Server: srv,
		Router: rtr,
		Cfg:    cfg,
	}
}

func (s *Server) Start() error {
	slog.Info("Starting HTTP server",
		"address", s.Cfg.Server.PortStr(),
		"read_timeout", s.Cfg.Server.ReadTimeout,
		"write_timeout", s.Cfg.Server.WriteTimeout,
		"idle_timeout", s.Cfg.Server.IdleTimeout,
	)

	return s.Server.ListenAndServe()
}
