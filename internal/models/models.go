package models

import "time"

type User struct {
	Nickname     string
	Email        string
	ID           string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
