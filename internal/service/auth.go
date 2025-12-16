package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	authRepository AuthRepository
	secretKey      string
	jwtExpiry      time.Duration
}

func NewAuthService(repository *postgres.Repository, secretKey string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		authRepository: repository,
		secretKey:      secretKey,
		jwtExpiry:      jwtExpiry,
	}
}

func (s AuthService) SignUp(ctx context.Context, nickname, email, password string) (*models.User, error) {
	const pp = "internal/service/auth.go/SignUp"

	slog.Info("Start creating User",
		"Path", pp,
	)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password",
			"Path", pp,
			"Error", err,
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
	const pp = "internal/service/auth.go/SignIn"

	slog.Info("Start user auth",
		"Path", pp,
	)

	user, err := s.authRepository.FindByEmail(ctx, email)
	if err != nil {
		slog.Warn("Failed to find email",
			"Path", pp,
			"Email", email,
			"Error", err,
		)
		return "", apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		slog.Warn("Failed to auth user",
			"Path", pp,
			"Email", email,
			"Error", err,
		)
		return "", apperrors.ErrInvalidCredentials
	}

	jwt, err := s.GenerateJWT(user)
	if err != nil {
		slog.Error("Failed to generate JWT",
			"Path", pp,
			"Error", err,
		)
		return "", err
	}

	slog.Info("JWt for user successfully created",
		"Path", pp,
		"Nickname", user.Nickname,
		"Email", user.Email,
	)
	return jwt, nil
}

type Claims struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

func (s AuthService) GenerateJWT(user *models.User) (string, error) {
	const pp = "internal/service/auth.go/GenerateJWT"

	claims := Claims{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtStr, err := jwtToken.SignedString([]byte(s.secretKey))
	if err != nil {
		slog.Error("Failed to sign token",
			"Path", pp,
			"Error", err,
		)
	}

	return jwtStr, nil
}
