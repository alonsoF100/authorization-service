package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
)

/*
pattern: /auth/register
method: POST
info: JSON in request body

succeed:

	-status code: 201 created
	-response body: JSON represented created user

failed:

	-status code: 400 bad request, 409 conflict, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	const pp = "internal/transport/http/handlers/auth.go/SignUp"
	
	var req dto.SignUpRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(ErrFailedToDecode))
		slog.Warn("Failed to decode JSON",
			"path", pp,
			"error", err,
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(ErrFailedToValidate))
		slog.Warn("Failed to validate request",
			"path", pp,
			"error", err,
		)
		return
	}

	slog.Debug("Data transfered to service layer",
		"nickname", req.Nickname,
		"email", req.Email,
		"passwordLength", len(req.Password),
	)
	user, err := h.userService.CreateUser(
		ctx,
		req.Nickname,
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserExist) {
			WriteJSON(w, http.StatusConflict, dto.NewErrorResponse(apperrors.ErrUserExist))
			slog.Debug("User with this nickname already exist",
				"username", req.Nickname,
				"error", err,
			)
			return
		}
		if errors.Is(err, apperrors.ErrEmailExist) {
			WriteJSON(w, http.StatusConflict, dto.NewErrorResponse(apperrors.ErrEmailExist))
			slog.Debug("User with this email already exist",
				"email", req.Email,
				"error", err,
			)
			return
		}
		WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(ErrServer))
		slog.Debug("Server error",
			"email", req.Email,
			"error", err,
		)
		return
	}

	WriteJSON(w, http.StatusCreated, dto.NewSignUpResponse(user))
}

/*
pattern: /auth/login
method: POST
info: JSON in request body

succeed:

	-status code: 200 ok
	-response body: JSON represented JWT token

failed:

	-status code: 400 bad request, 401 unauthorized, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	const pp = "internal/transport/http/handlers/user.go/SignIn"
}
