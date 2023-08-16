package create

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
)

type LocalKey struct {
	authenticityRepository repository.AuthenticityRepository
	localKeysRepository    repository.LocalKeysRepository
}

func NewLocalKey(authenticityRepository repository.AuthenticityRepository, localKeysRepository repository.LocalKeysRepository) *LocalKey {
	return &LocalKey{authenticityRepository: authenticityRepository, localKeysRepository: localKeysRepository}
}

func (k LocalKey) Create(ctx context.Context, keyType key.KeyType) (domain.LocalKey, error) {

	localKey, err := k.authenticityRepository.CreateLocalKey(keyType)
	if err != nil {
		return domain.LocalKey{}, err
	}

	err = k.localKeysRepository.SaveKey(ctx, localKey)
	if err != nil {
		return domain.LocalKey{}, err
	}
	return localKey, nil
}
