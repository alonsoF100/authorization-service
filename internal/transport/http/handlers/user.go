package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/alonsoF100/authorization-service/internal/transport/http/help"
	"github.com/alonsoF100/authorization-service/internal/transport/http/middleware"
)

/*
pattern: /api/me
method: GET
info: barer token from header

succeed:

	-status code: 200 ok
	-response body: JSON represented claims

failed:

	-status code: 401 unauthorized, 404 not found, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/user.go/GetMe"

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("User claims not found in context",
			slog.String("op", op))
		help.WriteJSON(w, http.StatusUnauthorized, dto.NewErrorResponse(apperrors.ErrUnauthorized))
		return
	}
	ctx := r.Context()

	user, err := h.UserService.GetUser(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFoundByID) {
			help.WriteJSON(w, http.StatusUnauthorized, dto.NewErrorResponse(apperrors.ErrInvalidCredentials))
			slog.Debug("Authentication failed",
				slog.String("op", op),
				slog.String("user_id", claims.ID),
				slog.String("error", err.Error()),
			)
			return
		}

		help.WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(apperrors.ErrServer))
		slog.Debug("Intenal server error",
			slog.String("op", op),
			slog.String("user_id", claims.ID),
			slog.String("error", err.Error()),
		)
		return
	}

	help.WriteJSON(w, http.StatusOK, dto.NewGetMeResponse(user))
}
