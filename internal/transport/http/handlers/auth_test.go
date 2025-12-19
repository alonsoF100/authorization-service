package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/alonsoF100/authorization-service/internal/transport/http/dto"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewAuthServiceMock(mc)
	mockValidator := validator.New()

	h := handlers.Handler{
		AuthService: mockService,
		UserService: nil,
		Validator:   mockValidator,
	}

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func()
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "bad JSON - missing coma",
			requestBody:    `{"name": "alex" "age": 10}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "emty JSON",
			requestBody:    ``,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "invalid JSON structure",
			requestBody:    `{name: "alex"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "failed validation - empty nickname",
			requestBody:    `{"nickname": "", "email": "test@test.com", "password": "password123"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToValidate,
		},
		{
			name:           "failed validation - invalid email",
			requestBody:    `{"nickname": "user", "email": "invalid-email", "password": "password123"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToValidate,
		},
		{
			name:           "failed validation - short password",
			requestBody:    `{"nickname": "user", "email": "test@test.com", "password": "123"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToValidate,
		},
		{
			name:        "user already exists by nickname",
			requestBody: `{"nickname": "existing", "email": "test@test.com", "password": "password123"}`,
			setupMocks: func() {
				mockService.SignUpMock.Expect(context.Background(), "existing", "test@test.com", "password123").Return(nil, apperrors.ErrUserExist)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  apperrors.ErrUserExist,
		},
		{
			name:        "user already exists by nickname",
			requestBody: `{"nickname": "alonso", "email": "existing@test.com", "password": "password123"}`,
			setupMocks: func() {
				mockService.SignUpMock.Expect(context.Background(), "alonso", "existing@test.com", "password123").Return(nil, apperrors.ErrEmailExist)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  apperrors.ErrEmailExist,
		},
		{
			name:        "user already exists by nickname",
			requestBody: `{"nickname": "alonso", "email": "alonso@test.com", "password": "password123"}`,
			setupMocks: func() {
				mockService.SignUpMock.Expect(context.Background(), "alonso", "alonso@test.com", "password123").Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  apperrors.ErrServer,
		},
		{
			name:        "user successfully created",
			requestBody: `{"nickname": "alonso", "email": "alonso@test.com", "password": "password123"}`,
			setupMocks: func() {
				mockService.SignUpMock.Expect(context.Background(), "alonso", "alonso@test.com", "password123").Return(&models.User{
					Nickname:     "alonso",
					Email:        "alonso@test.com",
					ID:           "some_uuid",
					PasswordHash: "password_hash",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			h.SignUp(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != nil {
				var errorResp dto.ErrorResponse

				err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				require.Equal(t, tt.expectedError.Error(), errorResp.Error)
				assert.NotEmpty(t, errorResp.TimeStamp)
			}
		})
	}
}

func TestSignUpSuccesses(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewAuthServiceMock(mc)
	mockValidator := validator.New()

	h := handlers.Handler{
		AuthService: mockService,
		UserService: nil,
		Validator:   mockValidator,
	}
	nickname := "alonsoF100"
	email := "alonso@mail.ru"
	password := "alonso_the_great"

	expectedUser := &models.User{
		Nickname:     nickname,
		Email:        email,
		ID:           uuid.New().String(),
		PasswordHash: "some_password_hash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockService.SignUpMock.Expect(context.Background(), nickname, email, password).Return(expectedUser, nil)

	requestBody := `{"nickname": "alonsoF100", "email": "alonso@mail.ru", "password": "alonso_the_great"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	h.SignUp(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var successResp dto.SignUpResponse
	err := json.Unmarshal(rr.Body.Bytes(), &successResp)
	require.NoError(t, err)

	assert.Equal(t, expectedUser.ID, successResp.UserID)
	assert.Equal(t, expectedUser.Nickname, successResp.Nickname)
	assert.Equal(t, expectedUser.Email, successResp.Email)
}

func TestSignIn(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewAuthServiceMock(mc)
	mockValidator := validator.New()

	h := handlers.Handler{
		AuthService: mockService,
		UserService: nil,
		Validator:   mockValidator,
	}

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func()
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "bad JSON - missing coma",
			requestBody:    `{"name": "alex" "age": 10}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "emty JSON",
			requestBody:    ``,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "invalid JSON structure",
			requestBody:    `{name: "alex"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToDecode,
		},
		{
			name:           "failed validation - invalid email",
			requestBody:    `{"email": "invalid-email", "password": "password123"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToValidate,
		},
		{
			name:           "failed validation - short password",
			requestBody:    `{"email": "test@test.com", "password": "123"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  apperrors.ErrFailedToValidate,
		},
		{
			name:        "invalid email or password",
			requestBody: `{"email": "alonso@mail.gaz", "password": "alonso_the_week"}`,
			setupMocks: func() {
				mockService.SignInMock.Expect(context.Background(), "alonso@mail.gaz", "alonso_the_week").Return("", apperrors.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  apperrors.ErrInvalidCredentials,
		},
		{
			name:        "server error",
			requestBody: `{"email": "alonso@mail.gaz", "password": "alonso_the_week"}`,
			setupMocks: func() {
				mockService.SignInMock.Expect(context.Background(), "alonso@mail.gaz", "alonso_the_week").Return("", errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  apperrors.ErrServer,
		},
		{
			name:        "success",
			requestBody: `{"email": "alonso@mail.ru", "password": "alonso_the_great"}`,
			setupMocks: func() {
				mockService.SignInMock.Expect(context.Background(), "alonso@mail.ru", "alonso_the_great").Return("oh_yes_JWT", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			h.SignIn(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != nil {
				var errorResp dto.ErrorResponse

				err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				require.Equal(t, tt.expectedError.Error(), errorResp.Error)
				assert.NotEmpty(t, errorResp.TimeStamp)
			}
		})
	}
}

func TestSignInSuccesses(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewAuthServiceMock(mc)
	mockValidator := validator.New()

	h := handlers.Handler{
		AuthService: mockService,
		UserService: nil,
		Validator:   mockValidator,
	}

	email := "alonso@mail.ru"
	password := "alonso_the_great"
	expectedToken := "valid token"
	mockService.SignInMock.Expect(context.Background(), email, password).Return(expectedToken, nil)

	requestBody := `{"email": "alonso@mail.ru", "password": "alonso_the_great"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	h.SignIn(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var successResp dto.SignInResponse
	err := json.Unmarshal(rr.Body.Bytes(), &successResp)
	require.NoError(t, err)

	require.Equal(t, expectedToken, successResp.JWT)
	require.Equal(t, "Bearer", successResp.Type)
}
