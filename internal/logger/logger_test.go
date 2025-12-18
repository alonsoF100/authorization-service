package logger_test

import (
	"log/slog"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel slog.Level
	}{
		{
			name:          "debug",
			level:         "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info",
			level:         "info",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn",
			level:         "warn",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error",
			level:         "error",
			expectedLevel: slog.LevelError,
		},
		{
			name:          "debug default",
			level:         "safdafasfafar",
			expectedLevel: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLevel, logger.ParseLevel(tt.level))
		})
	}
}

func TestSetup(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{"JSON debug", &config.Config{Logger: config.LoggerConfig{JSON: true, Level: "debug"}}},
		{"JSON info", &config.Config{Logger: config.LoggerConfig{JSON: true, Level: "info"}}},
		{"Text warn", &config.Config{Logger: config.LoggerConfig{JSON: false, Level: "warn"}}},
		{"Text error", &config.Config{Logger: config.LoggerConfig{JSON: false, Level: "error"}}},
		{"Default on invalid", &config.Config{Logger: config.LoggerConfig{JSON: true, Level: "invalid"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.Setup(tt.config)
			require.NotNil(t, log)
			require.Equal(t, log, slog.Default())
		})
	}
}
