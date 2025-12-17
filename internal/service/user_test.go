package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/service"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUserSuccess(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewUserRepositoryMock(mc)

	ctx := context.Background()
	userID := uuid.New().String()
	expectedUser := &models.User{
		ID:       userID,
		Nickname: "John Doe",
		Email:    "john@example.com",
	}

	mockRepo.FindByIDMock.Expect(ctx, userID).Return(expectedUser, nil)

	userService := service.NewUserService(mockRepo)

	user, err := userService.GetUser(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, expectedUser, user)
}

func TestGetUserDatabaseError(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewUserRepositoryMock(mc)

	ctx := context.Background()
	userID := uuid.New().String()
	someErr := errors.New("database error")

	mockRepo.FindByIDMock.Expect(ctx, userID).Return(nil, someErr)

	userService := service.NewUserService(mockRepo)

	user, err := userService.GetUser(ctx, userID)

	require.Error(t, err)
	require.True(t, errors.Is(err, someErr))
	require.Nil(t, user)
}

func TestGetUserNotFound(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := service.NewUserRepositoryMock(mc)

	ctx := context.Background()
	userID := uuid.New().String()
	someErr := apperrors.ErrUserNotFoundByID

	mockRepo.FindByIDMock.Expect(ctx, userID).Return(nil, nil)

	userService := service.NewUserService(mockRepo)

	user, err := userService.GetUser(ctx, userID)

	require.Error(t, err)
	require.Equal(t, someErr, err)
	require.Nil(t, user)
}
