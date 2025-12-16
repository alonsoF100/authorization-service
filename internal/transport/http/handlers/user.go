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
pattern: /auth/users
method: POST
info: JSON in request body

succeed:

	-status code: 201 created
	-response body: JSON represented created user

failed:

	-status code: 400 bad request, 409 conflict, 500 internal srv error
	-response body: JSON with err + time
*/
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const pp = "internal/transport/http/handlers/user.go/CreateUser"
	var req dto.CreateUserRequest
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

	WriteJSON(w, http.StatusCreated, dto.NewCreateUserResponse(user))
}
