package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r Repository) CreateUser(ctx context.Context, userDB *models.User) (*models.User, error) {
	const op = "repository/postgres/auth.go/CreateUser"

	const query = `
	INSERT INTO users (id, nickname, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING nickname, email, id, created_at 
	`

	slog.Debug("Query data",
		slog.String("op", op),
		slog.String("query_row", query),
		slog.String("id", userDB.ID),
		slog.String("nickname", userDB.Nickname),
		slog.String("email", userDB.Email),
		slog.Int("password_length", len(userDB.PasswordHash)),
		slog.Time("created_at", userDB.CreatedAt),
		slog.Time("updated_at", userDB.UpdatedAt),
	)
	var user models.User
	err := r.pool.QueryRow(
		ctx,
		query,
		userDB.ID,
		userDB.Nickname,
		userDB.Email,
		userDB.PasswordHash,
		userDB.CreatedAt,
		userDB.UpdatedAt,
	).Scan(
		&user.Nickname,
		&user.Email,
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "unique_email":
				slog.Debug("Email already exists",
					slog.String("op", op),
					slog.String("email", userDB.Email),
					slog.String("constraint", pgErr.ConstraintName),
				)
				return nil, apperrors.ErrEmailExist

			case "unique_nickname":
				slog.Debug("Nickname already exists",
					slog.String("op", op),
					slog.String("nickname", userDB.Nickname),
					slog.String("constraint", pgErr.ConstraintName),
				)
				return nil, apperrors.ErrUserExist
			}
		}

		slog.Error("Failed to create user",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slog.Debug("User created succsessfully",
		slog.String("op", op),
		slog.String("nickname", user.Nickname),
		slog.String("email", user.Email),
		slog.String("id", user.ID),
		slog.Time("created_at", user.CreatedAt),
	)

	return &user, nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "repository/postgres/auth.go/FindByEmail"

	const query = `
	SELECT id, email, nickname, password FROM users 
	WHERE email = $1
	`

	slog.Debug("Query data",
		slog.String("op", op),
		slog.String("query_row", query),
		slog.String("email", email),
	)
	var user models.User
	err := r.pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("User not found by email",
				slog.String("op", op),
				slog.String("email", email),
			)
			return nil, nil
		}

		slog.Error("Database error",
			slog.String("op", op),
			slog.String("email", email),
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
