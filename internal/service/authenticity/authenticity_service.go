package authenticity

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/authenticity/request"
	"context"
	"errors"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
)

type AuthenticityService struct {
	authenticityRepository repository.AuthenticityRepository
}

func NewAuthenticityService(authenticityRepository repository.AuthenticityRepository) *AuthenticityService {
	return &AuthenticityService{authenticityRepository: authenticityRepository}
}

var ErrKeyTypeNotSupported = errors.New("key type not supported for signing")

func (s AuthenticityService) Sign(ctx context.Context, request request.SignRequest) (string, []byte, error) {
	switch request.KeySource() {
	case domain.LOCAL_KEY:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithLocalKeyEns(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}

	case domain.MANAGED_KEY:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKeyEns(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}

	case domain.LOCAL_CERTIFICATE:
		return "", nil, nil
	case domain.MANAGED_CERTIFICATE:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKeyEns(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record.Retrieve(), nil
		}
	}

	return "", nil, ErrKeyTypeNotSupported
}
