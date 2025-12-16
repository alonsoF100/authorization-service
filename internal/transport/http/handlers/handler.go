package handlers

import (
	"context"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	SignUp(ctx context.Context, nickname, email, password string) (*models.User, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	authService AuthService
	validator   *validator.Validate
}

func New(authService AuthService) *Handler {
	return &Handler{
		authService: authService,
		validator:   validator.New(),
	}
}
