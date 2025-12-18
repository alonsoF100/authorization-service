package server_test

import (
	"testing"
	"time"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/transport/http/handlers"
	"github.com/alonsoF100/authorization-service/internal/transport/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	mockHandler := &handlers.Handler{}

	srv := server.New(cfg, mockHandler)

	require.NotNil(t, srv)
	assert.NotNil(t, srv.Server)
	assert.NotNil(t, srv.Router)
	assert.Equal(t, cfg, srv.Cfg)
	assert.Equal(t, ":8080", srv.Server.Addr)
	assert.Equal(t, 5*time.Second, srv.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, srv.Server.WriteTimeout)
	assert.Equal(t, 60*time.Second, srv.Server.IdleTimeout)
}
