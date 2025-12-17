package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r Repository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	const op = "repository/postgres/user.go/FindByID"

	const query = `
	SELECT id, nickname, email FROM users 
	WHERE id = $1
	`

	slog.Debug("Query data",
		slog.String("op", op),
		slog.String("query_row", query),
		slog.String("id", userID),
	)

	var user models.User
	err := r.pool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("User not found by id",
				slog.String("op", op),
				slog.String("user_id", userID),
			)
			return nil, nil
		}

		slog.Error("Database error",
			slog.String("op", op),
			slog.String("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slog.Debug("User was succsessfully founded",
		slog.String("op", op),
		slog.String("nickname", user.Nickname),
		slog.String("email", user.Email),
		slog.String("id", user.ID),
	)

	return &user, nil
}
