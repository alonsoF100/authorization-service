package apperrors

import "errors"

var (
	ErrUserExist          = errors.New("user with this nickname already exists")
	ErrEmailExist         = errors.New("user with this email already exists")
	ErrUserNotFoundByID   = errors.New("failed to find user by id")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUnauthorized       = errors.New("user authorized")
	ErrFailedToDecode     = errors.New("failed to decode JSON")
	ErrFailedToValidate   = errors.New("failed to validate request")
	ErrServer             = errors.New("damn, the server gaz up for nothing")
)
