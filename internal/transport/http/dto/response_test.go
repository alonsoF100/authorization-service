package dto_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/stretchr/testify/require"
)

func TestNewErrorResponse(t *testing.T) {
	someErr := errors.New("go go alonso go")
	response := dto.NewErrorResponse(someErr)

	require.Equal(t, someErr.Error(), response.Error)
	require.WithinDuration(t, time.Now(), response.TimeStamp, time.Millisecond)
}

func TestNewSignUpResponse(t *testing.T) {
	user := &models.User{
		Nickname:  "alonso",
		Email:     "alonso@goat.com",
		ID:        "33593c38-2a7a-4d94-b802-ed132a8fd4db",
		CreatedAt: time.Now(),
	}

	response := dto.NewSignUpResponse(user)

	require.Equal(t, user.Nickname, response.Nickname)
	require.Equal(t, user.Email, response.Email)
	require.Equal(t, user.ID, response.UserID)
	require.Equal(t, user.CreatedAt, response.CreatedAt)
}

func TestNewSignInResponse(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMzNTkzYzM4LTJhN2EtNGQ5NC1iODAyLWVkMTMyYThmZDRkYiIsImVtYWlsIjoiYWxvbnNvQG1haWwucnUiLCJuaWNrbmFtZSI6ImFsb25zb0YxMDAiLCJzdWIiOiIzMzU5M2MzOC0yYTdhLTRkOTQtYjgwMi1lZDEzMmE4ZmQ0ZGIiLCJleHAiOjE3NjYxNTMzNTMsImlhdCI6MTc2NjA2Njk1M30.13ulCHCkeDuQbz8IVnSRQZq9Ga-CvwcD6vQYvDH69yo"

	response := dto.NewSignInResponse(token)

	require.Equal(t, token, response.JWT)
	require.Equal(t, "Bearer", response.Type)
}

func TestNewGetMeResponse(t *testing.T) {
	user := &models.User{
		Nickname: "alonso",
		Email:    "alonso@goat.com",
		ID:       "33593c38-2a7a-4d94-b802-ed132a8fd4db",
	}

	response := dto.NewGetMeResponse(user)

	require.Equal(t, user.Nickname, response.Nickname)
	require.Equal(t, user.Email, response.Email)
	require.Equal(t, user.ID, response.ID)
}
