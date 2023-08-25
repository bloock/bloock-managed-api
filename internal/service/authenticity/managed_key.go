package authenticity

import (
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/authenticity/request"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
)

type ManagedKey struct {
	authenticityRepository repository.AuthenticityRepository
}

func NewManagedKey(authenticityRepository repository.AuthenticityRepository) *ManagedKey {
	return &ManagedKey{authenticityRepository: authenticityRepository}
}

func (k ManagedKey) Create(request request.CreateManagedKeyRequest) (key.ManagedKey, error) {
	managedKey, err := k.authenticityRepository.CreateManagedKey(
		request.Name(),
		request.KeyType(),
		request.Expiration(),
		request.Level(),
	)
	if err != nil {
		return key.ManagedKey{}, err
	}

	return managedKey, nil
}
