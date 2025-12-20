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
	Logger *slog.Logger
}

func New(cfg *config.Config, handlers *handlers.Handler, logger *slog.Logger) *Server {
	rtr := router.New(handlers).Setup()

	stdLogger := slog.NewLogLogger(logger.Handler(), slog.LevelError)

	srv := &http.Server{
		Addr:         cfg.Server.PortStr(),
		Handler:      rtr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorLog:     stdLogger,
	}

	return &Server{
		Server: srv,
		Router: rtr,
		Cfg:    cfg,
		Logger: logger,
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
