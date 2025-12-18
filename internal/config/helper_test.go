package config_test

import (
	"testing"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConStr(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.DatabaseConfig
		expectedStr string
	}{
		{
			name: "full config",
			config: &config.DatabaseConfig{
				User:     "postgres",
				Password: "postgres",
				Host:     "localhost",
				Port:     "5432",
				Name:     "auth",
				SSLMode:  "disable",
			},
			expectedStr: "postgresql://postgres:postgres@localhost:5432/auth?sslmode=disable",
		},
		{
			name: "with ssl require",
			config: &config.DatabaseConfig{
				User:     "admin",
				Password: "secret123",
				Host:     "db.example.com",
				Port:     "5432",
				Name:     "production",
				SSLMode:  "require",
			},
			expectedStr: "postgresql://admin:secret123@db.example.com:5432/production?sslmode=require",
		},
		{
			name: "with special characters in password",
			config: &config.DatabaseConfig{
				User:     "user",
				Password: "pass@word#123",
				Host:     "localhost",
				Port:     "5432",
				Name:     "test",
				SSLMode:  "disable",
			},
			expectedStr: "postgresql://user:pass@word#123@localhost:5432/test?sslmode=disable",
		},
		{
			name: "empty values",
			config: &config.DatabaseConfig{
				User:     "",
				Password: "",
				Host:     "",
				Port:     "",
				Name:     "",
				SSLMode:  "",
			},
			expectedStr: "postgresql://:@:/?sslmode=",
		},
		{
			name: "different port",
			config: &config.DatabaseConfig{
				User:     "user",
				Password: "pass",
				Host:     "127.0.0.1",
				Port:     "5433",
				Name:     "test",
				SSLMode:  "disable",
			},
			expectedStr: "postgresql://user:pass@127.0.0.1:5433/test?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ConStr()
			assert.Equal(t, tt.expectedStr, result)
		})
	}
}

func TestServerConfig_PortStr(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ServerConfig
		expected string
	}{
		{
			name:     "standard port",
			config:   &config.ServerConfig{Port: 8080},
			expected: ":8080",
		},
		{
			name:     "port 80",
			config:   &config.ServerConfig{Port: 80},
			expected: ":80",
		},
		{
			name:     "port 443",
			config:   &config.ServerConfig{Port: 443},
			expected: ":443",
		},
		{
			name:     "port 0",
			config:   &config.ServerConfig{Port: 0},
			expected: ":0",
		},
		{
			name:     "high port",
			config:   &config.ServerConfig{Port: 65535},
			expected: ":65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.PortStr()
			assert.Equal(t, tt.expected, result)
		})
	}
}
