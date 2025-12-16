package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/repository/postgres"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type AuthService struct {
	authRepository AuthRepository
}

func NewAuthService(repository *postgres.Repository) *AuthService {
	return &AuthService{
		authRepository: repository,
	}
}

func (s AuthService) SignUp(ctx context.Context, nickname, email, password string) (*models.User, error) {
	const pp = "internal/service/auth.go/SignUp"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password",
			"path", pp,
			"error", err,
		)
		return nil, err
	}

	userDB := &models.User{
		Nickname:     nickname,
		Email:        email,
		ID:           uuid.New().String(),
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	slog.Debug("User data transfered to db",
		"Path", pp,
		"Nickname", userDB.Nickname,
		"Email", userDB.Email,
	)
	user, err := s.authRepository.CreateUser(ctx, userDB)
	if err != nil {
		slog.Error("Failed to create user",
			"Path", pp,
			"Error", err,
		)
		return nil, err
	}

	slog.Info("User successfully created",
		"Path", pp,
		"UUID", user.ID,
		"Nickname", user.Nickname,
		"Email", user.Email,
	)
	return user, nil
}

func (s AuthService) SignIn(ctx context.Context, email, password string) (string, error) {

	return "", nil
}
