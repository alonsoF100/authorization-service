package postgres

import (
	"context"

	"github.com/alonsoF100/authorization-service/internal/models"
)

func (r Repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
