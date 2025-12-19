package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/service"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/router"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Basic(t *testing.T) {
	h := &handlers.Handler{
		AuthService: service.NewAuthService(nil, nil),
		UserService: service.NewUserService(nil),
		Validator:   nil,
	}

	rt := router.New(h)
	r := rt.Setup()

	testCases := []struct {
		method string
		path   string
		status int
	}{
		{"POST", "/auth/register", 400},
		{"POST", "/auth/login", 400},
		{"GET", "/api/me", 401},
		{"DELETE", "/api/me", 401},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.path == "/api/me" {
				req.Header.Set("Authorization", "Bearer token")
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.NotEqual(t, http.StatusNotFound, rr.Code,
				"path %s should be registered", tc.path)

			assert.GreaterOrEqual(t, rr.Code, tc.status,
				"status should be at least %d for %s", tc.status, tc.path)
		})
	}

	t.Run("non-existent path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/not-exist", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
