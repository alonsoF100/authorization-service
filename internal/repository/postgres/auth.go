package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alonsoF100/authorization-service/internal/apperrors"
	"github.com/alonsoF100/authorization-service/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r Repository) CreateUser(ctx context.Context, userDB *models.User) (*models.User, error) {
	const pp = "internal/repository/postgres/auth.go/CreateUser"

	const query = `
	INSERT INTO users (id, nickname, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING nickname, email, id, created_at 
	`

	slog.Debug("Query data",
		"Path", pp,
		"QueryRow", query,
		"ID", userDB.ID,
		"Nickname", userDB.Nickname,
		"Email", userDB.Email,
		"PasswordLength", len(userDB.PasswordHash),
		"CreatedAt", userDB.CreatedAt,
		"UpdatedAt", userDB.UpdatedAt,
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
				slog.Warn("Email already exists",
					"Path", pp,
					"Email", userDB.Email,
					"Constraint", pgErr.ConstraintName,
				)
				return nil, apperrors.ErrEmailExist

			case "unique_nickname":
				slog.Warn("Nickname already exists",
					"Path", pp,
					"Nickname", userDB.Nickname,
					"Constraint", pgErr.ConstraintName,
				)
				return nil, apperrors.ErrUserExist
			}
		}

		slog.Error("Failed to create user",
			"Path", pp,
			"Error", err,
			"ErrorType", fmt.Sprintf("%T", err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User created succsessfully",
		"Path", pp,
		"Nickname", user.Nickname,
		"Email", user.Email,
		"ID", user.ID,
		"CreatedAt", user.CreatedAt,
	)
	return &user, nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	const pp = "internal/repository/postgres/auth.go/FindByEmail"

	const query = `
	SELECT id, email, nickname, password FROM users 
	WHERE email = $1
	`

	slog.Debug("Query data",
		"Path", pp,
		"QueryRow", query,
		"Email", email,
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
		slog.Error("Failed to find user",
			"Path", pp,
			"Error", err,
			"ErrorType", fmt.Sprintf("%T", err),
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	slog.Info("User was succsessfully founded",
		"Path", pp,
		"Nickname", user.Nickname,
		"Email", user.Email,
		"ID", user.ID,
	)
	return &user, nil
}
