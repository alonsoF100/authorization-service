package dto

import (
	"time"

	"github.com/alonsoF100/authorization-service/internal/models"
)

type ErrorResponse struct {
	Error     string    `json:"error"`
	TimeStamp time.Time `json:"time_stamp"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:     err.Error(),
		TimeStamp: time.Now(),
	}
}

type SignUpResponse struct {
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSignUpResponse(user *models.User) SignUpResponse {
	return SignUpResponse{
		Nickname:  user.Nickname,
		Email:     user.Email,
		UserID:    user.ID,
		CreatedAt: user.CreatedAt,
	}
}
