package postgres

import (
	"context"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	const op = "repository/postgres/database.go/NewPool"

	poolConfig, err := pgxpool.ParseConfig(cfg.Database.ConStr())
	if err != nil {
		slog.Error("Failed to create pgx pool congig",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// TODO добавлять необходимые настройки pool а

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		slog.Error("Failed to create pgx pool",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		slog.Error("Failed to ping database",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	connConfig := poolConfig.ConnConfig
	db := stdlib.OpenDB(*connConfig)
	if err := goose.Up(db, cfg.Migration.Dir); err != nil {
		slog.Error("Failed to complete UP migrations",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return pool, nil
}
