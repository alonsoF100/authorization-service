package config

import "fmt"

func (cfg *DatabaseConfig) ConStr() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSlMode,
	)
}

func (cfg *ServerConfig) PortStr() string {
	return fmt.Sprintf(":%d", cfg.Port)
}
