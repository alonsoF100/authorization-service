package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/config"
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
	cfg            *config.Config
}

func NewAuthService(repository *postgres.Repository, cfg *config.Config) *AuthService {
	return &AuthService{
		authRepository: repository,
		cfg:            cfg,
	}
}

func (s AuthService) SignUp(ctx context.Context, nickname, email, password string) (*models.User, error) {
	const op = "service/auth.go/SignUp"

	slog.Debug("Start user registration",
		slog.String("op", op),
		slog.String("email", email),
	)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Registration failed: password hashing failed",
			slog.String("op", op),
			slog.String("email", email),
			slog.String("error", err.Error()),
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

	user, err := s.authRepository.CreateUser(ctx, userDB)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailExist) {
			slog.Info("Registration rejected: email already registered",
				slog.String("op", op),
				slog.String("email", email),
				slog.String("reason", "duplicate_email"),
			)
			return nil, apperrors.ErrEmailExist
		}
		if errors.Is(err, apperrors.ErrUserExist) {
			slog.Info("Registration rejected: nickname already taken",
				slog.String("op", op),
				slog.String("nickname", nickname),
				slog.String("reason", "duplicate_nickname"),
			)
			return nil, apperrors.ErrUserExist
		}

		slog.Error("Database error during registration",
			slog.String("op", op),
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slog.Info("Registration successful",
		slog.String("op", op),
		slog.String("user_id", user.ID),
		slog.String("email", email),
		slog.String("nickname", user.Nickname),
	)

	return user, nil
}

func (s AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	const op = "service/auth.go/SignIn"

	slog.Debug("Starting authentication",
		slog.String("op", op),
		slog.String("email", email),
	)

	user, err := s.authRepository.FindByEmail(ctx, email)
	if err != nil {
		slog.Error("Database error during authentication",
			slog.String("op", op),
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		slog.Info("Authentication failed: email not registered",
			slog.String("op", op),
			slog.String("email", email),
		)
		return "", apperrors.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		slog.Info("Authentication failed: invalid password",
			slog.String("op", op),
			slog.String("email", email),
			slog.String("user_id", user.ID),
		)
		return "", apperrors.ErrInvalidCredentials
	}

	jwt, err := s.GenerateJWT(user)
	if err != nil {
		slog.Error("Failed to generate JWT",
			slog.String("op", op),
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	slog.Info("Authentication successful",
		slog.String("op", op),
		slog.String("user_id", user.ID),
		slog.String("email", email),
		slog.String("nickname", user.Nickname),
	)

	return jwt, nil
}

func (s AuthService) GenerateJWT(user *models.User) (string, error) {
	const op = "service/auth.go/GenerateJWT"

	claims := models.Claims{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWT.Expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtStr, err := jwtToken.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		slog.Error("Failed to sign JWT token",
			slog.String("op", op),
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return jwtStr, nil
}

func (s AuthService) ValidateJWT(ctx context.Context, tokenString string) (*models.Claims, error) {
	const op = "service/auth.go/ValidateToken"

	var claims models.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Debug("Invalid signing method",
				slog.String("op", op),
				slog.String("method", t.Method.Alg()),
			)
			return nil, apperrors.ErrInvalidToken
		}
		return []byte(s.cfg.JWT.SecretKey), nil
	})
	if err != nil {
		slog.Debug("Token validation failed",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.ErrInvalidToken
	}

	if !token.Valid {
		slog.Debug("Token is invalid",
			slog.String("op", op),
		)
		return nil, apperrors.ErrInvalidToken
	}

	if time.Until(claims.ExpiresAt.Time) < 0 {
		slog.Debug("Token expired",
			slog.String("op", op),
			slog.String("user_id", claims.ID),
			slog.Time("expires_at", claims.ExpiresAt.Time),
		)
		return nil, apperrors.ErrTokenExpired
	}

	return &claims, nil
}
