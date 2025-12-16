package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsers, downCreateUsers)
}

func upCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE users (
			id UUID PRIMARY KEY,
			nickname VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(512) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMPNOT NULL
		);

		ALTER TABLE users 
            ADD CONSTRAINT unique_email UNIQUE (email);
        
        ALTER TABLE users 
            ADD CONSTRAINT unique_nickname UNIQUE (nickname);
	`)
	return err
}

func downCreateUsers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS users;")
	return err
}
