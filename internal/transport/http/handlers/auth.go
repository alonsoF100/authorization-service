package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/alonsoF100/authorization-service/internal/transport/http/help"
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
	const op = "handlers/auth.go/SignUp"

	var req dto.SignUpRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		help.WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(apperrors.ErrFailedToDecode))
		slog.Warn("Failed to decode JSON",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		help.WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(apperrors.ErrFailedToValidate))
		slog.Warn("Failed to validate request",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	user, err := h.AuthService.SignUp(
		ctx,
		req.Nickname,
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserExist) {
			help.WriteJSON(w, http.StatusConflict, dto.NewErrorResponse(apperrors.ErrUserExist))
			slog.Debug("User with this nickname already exist",
				slog.String("op", op),
				slog.String("nickname", req.Nickname),
				slog.String("error", err.Error()),
			)
			return
		}

		if errors.Is(err, apperrors.ErrEmailExist) {
			help.WriteJSON(w, http.StatusConflict, dto.NewErrorResponse(apperrors.ErrEmailExist))
			slog.Debug("User with this email already exist",
				slog.String("op", op),
				slog.String("email", req.Email),
				slog.String("error", err.Error()),
			)
			return
		}

		help.WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(apperrors.ErrServer))
		slog.Debug("Intenal server error",
			slog.String("op", op),
			slog.String("email", req.Email),
			slog.String("nickname", req.Nickname),
			slog.String("error", err.Error()),
		)
		return
	}

	help.WriteJSON(w, http.StatusCreated, dto.NewSignUpResponse(user))
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
	const op = "handlers/user.go/SignIn"

	var req dto.SignInRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		help.WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(apperrors.ErrFailedToDecode))
		slog.Warn("Failed to decode JSON",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		help.WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(apperrors.ErrFailedToValidate))
		slog.Warn("Failed to validate request",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	token, err := h.AuthService.SignIn(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			help.WriteJSON(w, http.StatusUnauthorized, dto.NewErrorResponse(apperrors.ErrInvalidCredentials))
			slog.Debug("Authentication failed",
				slog.String("op", op),
				slog.String("email", req.Email),
				slog.String("error", err.Error()),
			)
			return
		}

		help.WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(apperrors.ErrServer))
		slog.Debug("Intenal server error",
			slog.String("op", op),
			slog.String("email", req.Email),
			slog.String("error", err.Error()),
		)
		return
	}

	help.WriteJSON(w, http.StatusOK, dto.NewSignInResponse(token))
}
