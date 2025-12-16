package handlers

import (
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

	-status code: 401 unauthorized, 500 internal server error
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

	// TODO get valid data from database

	slog.Info("User info requested",
		slog.String("op", op),
		slog.String("user_id", claims.ID),
		slog.String("email", claims.Email),
		slog.String("nickname", claims.Nickname),
	)

	help.WriteJSON(w, http.StatusOK, dto.NewGetMeResponse(claims))
}
