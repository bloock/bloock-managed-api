package encryption

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/encryption/request"
	"context"
	"errors"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type EncryptionService struct {
	encryptionRepository repository.EncryptionRepository
}

func NewEncryptionService(encryptionRepository repository.EncryptionRepository) *EncryptionService {
	return &EncryptionService{
		encryptionRepository: encryptionRepository,
	}
}

var ErrKeyTypeNotSupported = errors.New("key type not supported for encrypting")

func (s EncryptionService) Encrypt(ctx context.Context, request request.EncryptRequest) (*record.Record, error) {
	switch request.KeySource() {
	case domain.LOCAL_KEY:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return record, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey())
			if err != nil {
				return nil, err
			}

			return record, nil
		}

	case domain.MANAGED_KEY:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		}

	case domain.LOCAL_CERTIFICATE:
		return nil, nil
	case domain.MANAGED_CERTIFICATE:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		}
	}

	return nil, ErrKeyTypeNotSupported
}
