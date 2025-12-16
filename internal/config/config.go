package config

import "time"

type Config struct {
	Server    ServerConfig     `mapstructure:"server"`
	Database  DatabaseConfig   `mapstructure:"database"`
	Logger    LoggerConfig     `mapstructure:"logger"`
	Migration MigrationsConfig `mapstructure:"migrations"`
	JWT       JWTConfig        `mapstructure:"jwt"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
	JSON  bool   `mapstructure:"json"`
}

type MigrationsConfig struct {
	Dir string `mapstructure:"dir"`
}

type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key"`
	Expiry    time.Duration `mapstructure:"expiry"`
}
