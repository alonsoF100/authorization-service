package handlers

import (
	"context"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	// TODO add methods
	CreateUser(ctx context.Context, nickname, email, password string) (*models.User, error)
}

type Handler struct {
	userService UserService
	validator   *validator.Validate
}

func New(userService UserService) *Handler {
	return &Handler{
		userService: userService,
		validator:   validator.New(),
	}
}
