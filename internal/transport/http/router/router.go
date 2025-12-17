package router

import (
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	handlers *handlers.Handler
}

func New(handlers *handlers.Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (rt Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", rt.handlers.SignUp)
		r.Post("/login", rt.handlers.SignIn)
	})

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Auth(rt.handlers.AuthService))

		r.Get("/me", rt.handlers.GetMe)
		r.Delete("/me", rt.handlers.DeleteMe)
	})

	return r
}
