package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/repository/postgres"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (*models.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(repository *postgres.Repository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (s UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	const op = "service/user.go/GetUser"

	slog.Debug("Start invalidation user data",
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		slog.Error("Database error during invalidation user data",
			slog.String("op", op),
			slog.String("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		slog.Info("Authentication failed: email not registered",
			slog.String("op", op),
			slog.String("user_id", userID),
		)
		return nil, apperrors.ErrUserNotFoundByID
	}

	slog.Info("User founded successfuly",
		slog.String("op", op),
		slog.String("user_id", user.ID),
		slog.String("email", user.Email),
		slog.String("nickname", user.Nickname),
	)

	return user, nil
}
