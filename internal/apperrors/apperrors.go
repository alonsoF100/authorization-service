package apperrors

import "errors"

var (
	ErrUserExist          = errors.New("user with this nickname already exists")
	ErrEmailExist         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
