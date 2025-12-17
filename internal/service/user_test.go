package service_test

import (
	"context"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/service"
	"github.com/alonsoF100/authorization-service/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSuccess(t *testing.T) {
	mc := minimock.NewController(t)
	mockRepo := mocks.NewUserRepositoryMock(mc)

	ctx := context.Background()
	userID := "123"
	expectedUser := &models.User{
		ID:       "123",
		Nickname: "John Doe",
		Email:    "john@example.com",
	}

	mockRepo.FindByIDMock.Expect(ctx, userID).Return(expectedUser, nil)

	userService := service.NewUserService(mockRepo)

	user, err := userService.GetUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
}
