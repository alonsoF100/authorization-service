package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	expectedCfg := config.Config{
		Server: config.ServerConfig{
			Port:         8080,
			ReadTimeout:  time.Duration(5) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second,
			IdleTimeout:  time.Duration(10) * time.Second,
		},
		Database: config.DatabaseConfig{
			Host:     "postgres",
			Port:     "5432",
			User:     "testuser",
			Password: "testpassword",
			Name:     "auth",
			SSLMode:  "disable",
		},
		Logger: config.LoggerConfig{
			Level: "info",
			JSON:  false,
		},
		Migration: config.MigrationsConfig{
			Dir: "migrations/postgres",
		},
		JWT: config.JWTConfig{
			Expiry:    time.Duration(24) * time.Hour,
			SecretKey: "test-secret-key",
		},
	}

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	envContent := `DB_USER=testuser
DB_PASSWORD=testpassword
SECRET_KEY=test-secret-key`

	require.NoError(t, os.WriteFile(".env", []byte(envContent), 0644))

	yamlContent := `server:
  port: 8080
  read_timeout: "5s"
  write_timeout: "10s"
  idle_timeout: "10s"

database:
  host: postgres
  port: "5432"
  name: auth
  ssl_mode: disable

logger:
  level: "info"
  json: false

migrations:
  dir: "migrations/postgres"

jwt:
  expiry: "24h"`

	require.NoError(t, os.WriteFile("config.yaml", []byte(yamlContent), 0644))

	cfg := config.Load()
	require.NotNil(t, cfg)

	require.Equal(t, expectedCfg.Server.Port, cfg.Server.Port)
	require.Equal(t, expectedCfg.Server.ReadTimeout, cfg.Server.ReadTimeout)
	require.Equal(t, expectedCfg.Server.WriteTimeout, cfg.Server.WriteTimeout)
	require.Equal(t, expectedCfg.Server.IdleTimeout, cfg.Server.IdleTimeout)

	require.Equal(t, expectedCfg.Database.Host, cfg.Database.Host)
	require.Equal(t, expectedCfg.Database.Port, cfg.Database.Port)
	require.Equal(t, expectedCfg.Database.Name, cfg.Database.Name)
	require.Equal(t, expectedCfg.Database.SSLMode, cfg.Database.SSLMode)
	require.Equal(t, expectedCfg.Database.User, cfg.Database.User)
	require.Equal(t, expectedCfg.Database.Password, cfg.Database.Password)

	require.Equal(t, expectedCfg.Logger.JSON, cfg.Logger.JSON)
	require.Equal(t, expectedCfg.Logger.Level, cfg.Logger.Level)

	require.Equal(t, expectedCfg.Migration.Dir, cfg.Migration.Dir)

	require.Equal(t, expectedCfg.JWT.Expiry, expectedCfg.JWT.Expiry)
	require.Equal(t, expectedCfg.JWT.SecretKey, cfg.JWT.SecretKey)
}
