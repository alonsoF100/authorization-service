package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Nickname     string
	Email        string
	ID           string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Claims struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}
