package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/service"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpSuccess(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	nickname := "alonsoF100"
	email := "alonso@yandex.ru"
	password := "alonso_the_great"

	mockRepo.CreateUserMock.Set(func(ctx context.Context, user *models.User) (up1 *models.User, err error) {
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.PasswordHash)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
		require.NotEqual(t, password, user.PasswordHash)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		require.NoError(t, err)

		return user, nil
	})

	authService := service.NewAuthService(mockRepo, nil)

	user, err := authService.SignUp(ctx, nickname, email, password)

	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, nickname, user.Nickname)
	require.Equal(t, email, user.Email)
	require.NotEmpty(t, user.ID)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	require.NoError(t, err)
}

func TestSignUpEmailAlreadyExist(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	nickname := "alonsoF100"
	email := "alonso@yandex.ru"
	password := "alonso_the_great"

	mockRepo.CreateUserMock.Set(func(ctx context.Context, user *models.User) (up1 *models.User, err error) {
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.PasswordHash)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
		require.NotEqual(t, password, user.PasswordHash)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		require.NoError(t, err)

		return nil, apperrors.ErrEmailExist
	})

	authService := service.NewAuthService(mockRepo, nil)

	user, err := authService.SignUp(ctx, nickname, email, password)

	require.Error(t, err)
	require.Nil(t, user)
	require.True(t, errors.Is(err, apperrors.ErrEmailExist))
}

func TestSignUpNicknameAlreadyExist(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	nickname := "alonsoF100"
	email := "alonso@yandex.ru"
	password := "alonso_the_great"

	mockRepo.CreateUserMock.Set(func(ctx context.Context, user *models.User) (up1 *models.User, err error) {
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.PasswordHash)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
		require.NotEqual(t, password, user.PasswordHash)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		require.NoError(t, err)

		return nil, apperrors.ErrUserExist
	})

	authService := service.NewAuthService(mockRepo, nil)

	user, err := authService.SignUp(ctx, nickname, email, password)

	require.Error(t, err)
	require.Nil(t, user)
	require.True(t, errors.Is(err, apperrors.ErrUserExist))
}

func TestSignUpDatabaseError(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	nickname := "alonsoF100"
	email := "alonso@yandex.ru"
	password := "alonso_the_great"
	someErr := errors.New("database error")

	mockRepo.CreateUserMock.Set(func(ctx context.Context, user *models.User) (up1 *models.User, err error) {
		require.Equal(t, nickname, user.Nickname)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.PasswordHash)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
		require.NotEqual(t, password, user.PasswordHash)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		require.NoError(t, err)

		return nil, someErr
	})

	authService := service.NewAuthService(mockRepo, nil)

	user, err := authService.SignUp(ctx, nickname, email, password)

	require.Error(t, err)
	require.Nil(t, user)
	require.True(t, errors.Is(err, someErr))
}

func TestSignInSuccess(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	email := "alonso@yandex.ru"
	password := "alonso_the_great"
	config := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "someSecret",
			Expiry:    time.Duration(24) * time.Hour,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	expectedUser := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Nickname:     "alonsoF100",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.FindByEmailMock.Expect(ctx, email).Return(expectedUser, nil)

	authService := service.NewAuthService(mockRepo, config)

	jwtT, err := authService.SignIn(ctx, email, password)

	require.NoError(t, err)
	require.NotEmpty(t, jwtT)

	var claims models.Claims

	token, err := jwt.ParseWithClaims(jwtT, &claims, func(t *jwt.Token) (any, error) {
		return []byte(config.JWT.SecretKey), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
	require.Equal(t, email, claims.Email)
	require.Equal(t, expectedUser.Nickname, claims.Nickname)
	require.Equal(t, expectedUser.ID, claims.ID)
}

func TestSignInDatabaseError(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	email := "alonso@yandex.ru"
	password := "alonso_the_great"
	someErr := errors.New("database error")
	config := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "someSecret",
			Expiry:    time.Duration(24) * time.Hour,
		},
	}

	mockRepo.FindByEmailMock.Expect(ctx, email).Return(nil, someErr)

	authService := service.NewAuthService(mockRepo, config)

	jwt, err := authService.SignIn(ctx, email, password)

	require.Error(t, err)
	require.Empty(t, jwt)
	require.True(t, errors.Is(err, someErr))
}

func TestSignInWrongEmail(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	email := "alonso@yandex.ru"
	password := "alonso_the_great"
	config := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "someSecret",
			Expiry:    time.Duration(24) * time.Hour,
		},
	}

	mockRepo.FindByEmailMock.Expect(ctx, email).Return(nil, nil)

	authService := service.NewAuthService(mockRepo, config)

	jwt, err := authService.SignIn(ctx, email, password)

	require.Error(t, err)
	require.Empty(t, jwt)
	require.True(t, errors.Is(err, apperrors.ErrInvalidCredentials))
}

func TestSignInWrongPassword(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewAuthRepositoryMock(mc)

	ctx := context.Background()
	email := "alonso@yandex.ru"
	password := "alonso_the_great"
	wrongPassword := "alonso_the_worst"
	config := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "someSecret",
			Expiry:    time.Duration(24) * time.Hour,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	expectedUser := &models.User{
		PasswordHash: string(hashedPassword),
	}

	mockRepo.FindByEmailMock.Expect(ctx, email).Return(expectedUser, nil)

	authService := service.NewAuthService(mockRepo, config)

	jwt, err := authService.SignIn(ctx, email, wrongPassword)

	require.Error(t, err)
	require.Empty(t, jwt)
	require.True(t, errors.Is(err, apperrors.ErrInvalidCredentials))
}

func TestValidateJWT(t *testing.T) {
	config := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "someSecretsomeSecretsomeSecretsomeSecretsomeSecretsomeSecretsomeSecretsomeSecret",
			Expiry:    time.Duration(24) * time.Hour,
		},
	}

	authService := service.NewAuthService(nil, config)

	tests := []struct {
		name      string
		token     string
		wantError bool
		errorType error
	}{
		{
			name:      "wrongMethod",
			token:     "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxNTE2MjM5MDIzfQ.y-_R9WFHPlN7lFbHXbflMbF0jm4i6Ow0ZNFOOcDfo3s8IkDyBtvXZ_kdD31LAzQiWhMflPK4gAXRkwA8jIg7Gw",
			wantError: true,
			errorType: apperrors.ErrInvalidToken,
		},
		{
			name:      "not a jwt",
			token:     "gaz gaz",
			wantError: true,
			errorType: apperrors.ErrInvalidToken,
		},
		{
			name:      "expired token",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMzNTkzYzM4LTJhN2EtNGQ5NC1iODAyLWVkMTMyYThmZDRkYiIsImVtYWlsIjoiYWxvbnNvQG1haWwucnUiLCJuaWNrbmFtZSI6ImFsb25zb0YxMDAiLCJzdWIiOiIzMzU5M2MzOC0yYTdhLTRkOTQtYjgwMi1lZDEzMmE4ZmQ0ZGIiLCJleHAiOjE3NjYwNjY0NTAsImlhdCI6MTc2NjA2NjQ1MH0.i1Ctjl4n7l2n2L3S4NGjQ1h9481xqh0z0E2rk2mUUA4",
			wantError: true,
			errorType: apperrors.ErrInvalidToken,
		},
		{
			name:      "wrong secret",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMzNTkzYzM4LTJhN2EtNGQ5NC1iODAyLWVkMTMyYThmZDRkYiIsImVtYWlsIjoiYWxvbnNvQG1haWwucnUiLCJuaWNrbmFtZSI6ImFsb25zb0YxMDAiLCJzdWIiOiIzMzU5M2MzOC0yYTdhLTRkOTQtYjgwMi1lZDEzMmE4ZmQ0ZGIiLCJleHAiOjE3NjYxNTA3NDYsImlhdCI6MTc2NjA2NDM0Nn0.M2XEijWjOCiN17ttRCwucfiMyPUlDmedlwZLN5R8pK8",
			wantError: true,
			errorType: apperrors.ErrInvalidToken,
		},
		{
			name:      "good token",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMzNTkzYzM4LTJhN2EtNGQ5NC1iODAyLWVkMTMyYThmZDRkYiIsImVtYWlsIjoiYWxvbnNvQG1haWwucnUiLCJuaWNrbmFtZSI6ImFsb25zb0YxMDAiLCJzdWIiOiIzMzU5M2MzOC0yYTdhLTRkOTQtYjgwMi1lZDEzMmE4ZmQ0ZGIiLCJleHAiOjE3NjYxNTA3NDYsImlhdCI6MTc2NjA2NDM0Nn0.bCl1wqh59UuBiSGGqH8bRo_PJ-6-U2BjhrgUDGzBHQM",
			wantError: false,
			errorType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateJWT(context.Background(), tt.token)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorType != nil {
					require.True(t, errors.Is(err, tt.errorType),
						"Expected error type %v, got %v", tt.errorType, err)
				}
				require.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
			}
		})
	}
}
