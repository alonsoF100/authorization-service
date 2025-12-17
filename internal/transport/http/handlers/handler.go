package handlers

import (
	"context"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	SignUp(ctx context.Context, nickname, email, password string) (*models.User, error)
	SignIn(ctx context.Context, email, password string) (string, error)
	ValidateJWT(ctx context.Context, tokenString string) (*models.Claims, error)
}

type UserService interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type Handler struct {
	AuthService AuthService
	UserService UserService
	validator   *validator.Validate
}

func New(authService AuthService, userService UserService) *Handler {
	return &Handler{
		AuthService: authService,
		UserService: userService,
		validator:   validator.New(),
	}
}
