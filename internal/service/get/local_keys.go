package get

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
)

type LocalKeys struct {
	localKeysRepository repository.LocalKeysRepository
}

func NewLocalKeys(localKeysRepository repository.LocalKeysRepository) *LocalKeys {
	return &LocalKeys{localKeysRepository: localKeysRepository}
}

func (s LocalKeys) Get(ctx context.Context) ([]domain.LocalKey, error) {
	return s.localKeysRepository.FindKeys(ctx)
}
