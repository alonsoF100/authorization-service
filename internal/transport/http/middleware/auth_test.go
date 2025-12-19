package middleware_test

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
	"github.com/alonsoF100/authorization-service/internal/transport/http/middleware"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "valid Bearer token",
			header: "Bearer token",
			want:   "token",
		},
		{
			name:   "empty header",
			header: "",
			want:   "",
		},
		{
			name:   "Bearer without token",
			header: "Bearer",
			want:   "",
		},
		{
			name:   "token without Bearer",
			header: "token",
			want:   "",
		},
		{
			name:   "too many parts",
			header: "Bearer token extra",
			want:   "",
		},
		{
			name:   "wrong prefix",
			header: "Basic token",
			want:   "",
		},
		{
			name:   "Bearer with JWT",
			header: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/any-url", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			got := middleware.ExtractToken(req)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAuth(t *testing.T) {
	mc := minimock.NewController(t)

	mockValidator := middleware.NewTokenValidatorMock(mc)

	testClaims := &models.Claims{
		ID:       uuid.New().String(),
		Nickname: "alonso",
		Email:    "alonso@mail.ru",
	}

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		setupMocks     func()
		expectedStatus int
		expectedError  error
		shouldCallNext bool
	}{
		{
			name: "successful authentication",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid_token")
			},
			setupMocks: func() {
				mockValidator.ValidateJWTMock.Expect(context.Background(), "valid_token").
					Return(testClaims, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
			shouldCallNext: true,
		},
		{
			name:           "missing authorization header",
			setupRequest:   func(req *http.Request) {},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  middleware.ErrNoAuthHeader,
			shouldCallNext: false,
		},
		{
			name: "empty authorization header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  middleware.ErrNoAuthHeader,
			shouldCallNext: false,
		},
		{
			name: "invalid token format - no Bearer",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "invalid_token")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  middleware.ErrNoAuthHeader,
			shouldCallNext: false,
		},
		{
			name: "invalid token format - Bearer without token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  middleware.ErrNoAuthHeader,
			shouldCallNext: false,
		},
		{
			name: "invalid token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid_token")
			},
			setupMocks: func() {
				mockValidator.ValidateJWTMock.Expect(context.Background(), "invalid_token").
					Return(nil, apperrors.ErrInvalidToken)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  apperrors.ErrInvalidToken,
			shouldCallNext: false,
		},
		{
			name: "validator returns unexpected error",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer token_with_error")
			},
			setupMocks: func() {
				mockValidator.ValidateJWTMock.Expect(context.Background(), "token_with_error").
					Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  apperrors.ErrInvalidToken,
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			nextCalled := false
			var capturedClaims interface{}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				claims := r.Context().Value(middleware.UserContextKey)
				capturedClaims = claims

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			authMiddleware := middleware.Auth(mockValidator)
			handler := authMiddleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			tt.setupRequest(req)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedStatus, rr.Code)
			require.Equal(t, tt.shouldCallNext, nextCalled)

			if tt.expectedError != nil {
				var errorResp dto.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedError.Error(), errorResp.Error)

				assert.NotEmpty(t, errorResp.TimeStamp)
			}

			if tt.shouldCallNext {
				require.NotNil(t, capturedClaims)
				require.Equal(t, testClaims, capturedClaims)
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	// Создаем тестовые claims
	testClaims := &models.Claims{
		ID:       uuid.New().String(),
		Nickname: "alonso",
		Email:    "alonso@gleb.gazz",
	}

	tests := []struct {
		name      string
		setupCtx  func() context.Context
		wantUser  *models.Claims
		wantFound bool
	}{
		{
			name: "user found in context",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), middleware.UserContextKey, testClaims)
			},
			wantUser:  testClaims,
			wantFound: true,
		},
		{
			name: "user not found - empty context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantUser:  nil,
			wantFound: false,
		},
		{
			name: "user not found - wrong key",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), "wrong_key", testClaims)
			},
			wantUser:  nil,
			wantFound: false,
		},
		{
			name: "user not found - nil value",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), middleware.UserContextKey, nil)
			},
			wantUser:  nil,
			wantFound: false,
		},
		{
			name: "user not found - wrong type in value",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), middleware.UserContextKey, "not_claims")
			},
			wantUser:  nil,
			wantFound: false,
		},
		{
			name: "user found with additional values in context",
			setupCtx: func() context.Context {
				ctx := context.WithValue(context.Background(), "other_key", "other_value")
				return context.WithValue(ctx, middleware.UserContextKey, testClaims)
			},
			wantUser:  testClaims,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()

			gotUser, gotFound := middleware.GetUserFromContext(ctx)

			require.Equal(t, tt.wantFound, gotFound)

			if tt.wantFound {
				require.NotNil(t, gotUser)
				assert.Equal(t, tt.wantUser.ID, gotUser.ID)
				assert.Equal(t, tt.wantUser.Nickname, gotUser.Nickname)
				assert.Equal(t, tt.wantUser.Email, gotUser.Email)
			} else {
				assert.Nil(t, gotUser)
			}
		})
	}
}
