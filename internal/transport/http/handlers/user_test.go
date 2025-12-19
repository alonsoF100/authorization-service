package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/middleware"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetMe(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewUserServiceMock(mc)

	h := handlers.Handler{
		UserService: mockService,
	}

	tests := []struct {
		name       string
		claims     *models.Claims
		mockSetup  func(ctx context.Context)
		wantStatus int
		wantError  string
	}{
		{
			name: "success",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.GetUserMock.Expect(ctx, "user123").
					Return(&models.User{
						ID:       "user123",
						Nickname: "testuser",
						Email:    "test@test.com",
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantError:  "",
		},
		{
			name:   "no claims in context",
			claims: nil,
			mockSetup: func(ctx context.Context) {
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  apperrors.ErrUnauthorized.Error(),
		},
		{
			name: "user not found",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.GetUserMock.Expect(ctx, "user123").
					Return(nil, apperrors.ErrUserNotFoundByID)
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  apperrors.ErrInvalidCredentials.Error(),
		},
		{
			name: "service error",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.GetUserMock.Expect(ctx, "user123").
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  apperrors.ErrServer.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.claims != nil {
				ctx = context.WithValue(ctx, middleware.UserContextKey, tt.claims)
			}

			tt.mockSetup(ctx)

			req := httptest.NewRequest("GET", "/api/me", nil)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			h.GetMe(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)

			if tt.wantError != "" {
				var resp dto.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, tt.wantError, resp.Error)
			}
		})
	}
}

func TestDeleteMe(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewUserServiceMock(mc)

	h := handlers.Handler{
		UserService: mockService,
	}

	tests := []struct {
		name       string
		claims     *models.Claims
		mockSetup  func(ctx context.Context)
		wantStatus int
		wantError  string
	}{
		{
			name: "success - user deleted",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.DeleteUserMock.Expect(ctx, "user123").
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
			wantError:  "",
		},
		{
			name:   "no claims in context",
			claims: nil,
			mockSetup: func(ctx context.Context) {
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  apperrors.ErrUnauthorized.Error(),
		},
		{
			name: "user not found",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.DeleteUserMock.Expect(ctx, "user123").
					Return(apperrors.ErrUserNotFoundByID)
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  apperrors.ErrInvalidCredentials.Error(),
		},
		{
			name: "service error",
			claims: &models.Claims{
				ID: "user123",
			},
			mockSetup: func(ctx context.Context) {
				mockService.DeleteUserMock.Expect(ctx, "user123").
					Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  apperrors.ErrServer.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.claims != nil {
				ctx = context.WithValue(ctx, middleware.UserContextKey, tt.claims)
			}

			tt.mockSetup(ctx)

			req := httptest.NewRequest("DELETE", "/api/me", nil)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			h.DeleteMe(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)

			if tt.wantError != "" {
				var resp dto.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, tt.wantError, resp.Error)
			}

			if tt.wantStatus == http.StatusNoContent {
				require.Equal(t, "null\n", rr.Body.String())
			}
		})
	}
}
