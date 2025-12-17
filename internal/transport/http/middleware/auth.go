package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/alonsoF100/authorization-service/internal/transport/http/help"
)

type TokenValidator interface {
	ValidateJWT(ctx context.Context, token string) (*models.Claims, error)
}

var (
	ErrNoAuthHeader = errors.New("Missing authorization header")
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(tokenValidator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware/auth.go/Auth"

			token := extractToken(r)
			if token == "" {
				slog.Info("Authentication failed: missing authorization header",
					slog.String("op", op),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
				help.WriteJSON(w, http.StatusUnauthorized, dto.NewErrorResponse(ErrNoAuthHeader))
				return
			}

			claims, err := tokenValidator.ValidateJWT(r.Context(), token)
			if err != nil {
				slog.Info("Authentication failed: invalid token",
					slog.String("op", op),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
				help.WriteJSON(w, http.StatusUnauthorized, dto.NewErrorResponse(apperrors.ErrInvalidToken))
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func GetUserFromContext(ctx context.Context) (*models.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*models.Claims)

	return claims, ok
}
