package repository

import (
	"bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/authenticity"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type BloockAuthenticityRepository struct {
	apikey             string
	keyClient          client.KeyClient
	authenticityClient client.AuthenticityClient
	recordClient       client.RecordClient
	log                zerolog.Logger
}

func NewBloockAuthenticityRepository(apikey string, keyClient client.KeyClient, authenticityClient client.AuthenticityClient, recordClient client.RecordClient, log zerolog.Logger) *BloockAuthenticityRepository {
	return &BloockAuthenticityRepository{apikey: apikey, keyClient: keyClient, authenticityClient: authenticityClient, recordClient: recordClient, log: log}
}

func (b BloockAuthenticityRepository) Sign(localKey key.LocalKey, keyType key.KeyType, commonName *string, data []byte) (record.Record, error) {
	var err error
	var rc record.Record
	if keyType == key.EcP256k {
		rc, err = b.recordClient.FromBytes(data).WithSigner(authenticity.NewEcdsaSigner(authenticity.SignerArgs{
			LocalKey:   &localKey,
			CommonName: commonName,
		})).Build()
	} else {
		rc, err = b.recordClient.FromBytes(data).WithSigner(authenticity.NewEnsSigner(authenticity.SignerArgs{
			LocalKey:   &localKey,
			CommonName: commonName,
		})).Build()
	}

	if err != nil {
		return record.Record{}, err
	}

	return rc, nil
}

func (b BloockAuthenticityRepository) Verify(record record.Record) (bool, error) {
	return b.authenticityClient.Verify(record)
}

func (b BloockAuthenticityRepository) CreateLocalKey(keyType key.KeyType) (domain.LocalKey, error) {
	localKey, err := b.keyClient.NewLocalKey(keyType)
	if err != nil {
		return domain.LocalKey{}, err
	}

	return *domain.NewLocalKey(localKey, keyType, uuid.New()), nil
}

func (b BloockAuthenticityRepository) CreateManagedKey(name string, keyType key.KeyType, expiration int, level key.KeyProtectionLevel) (key.ManagedKey, error) {
	return b.keyClient.NewManagedKey(key.ManagedKeyParams{
		Name:       name,
		Protection: level,
		KeyType:    keyType,
		Expiration: int64(expiration),
	})
}

func (b BloockAuthenticityRepository) LoadLocalKey(keyType key.KeyType, publicKey string, privateKey string) (key.LocalKey, error) {
	return b.keyClient.LoadLocalKey(keyType, publicKey, &privateKey)
}
