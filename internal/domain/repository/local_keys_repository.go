package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
	"github.com/google/uuid"
)

type LocalKeysRepository interface {
	SaveKey(ctx context.Context, localKey domain.LocalKey) error
	FindKeyByID(ctx context.Context, id uuid.UUID) (*domain.LocalKey, error)
	FindKeys(ctx context.Context) ([]domain.LocalKey, error)
}
