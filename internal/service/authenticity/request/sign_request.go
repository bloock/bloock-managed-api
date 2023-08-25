package request

import (
	"bloock-managed-api/internal/service"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type SignRequest struct {
	publicKey        string
	privateKey       *string
	keyID            uuid.UUID
	kty              key.KeyType
	keyType          service.KeyType
	data             []byte
	useEnsResolution bool
}

func (s SignRequest) PublicKey() string {
	return s.publicKey
}

func (s SignRequest) PrivateKey() *string {
	return s.privateKey
}

func (s SignRequest) KeyID() uuid.UUID {
	return s.keyID
}

func (s SignRequest) Kty() key.KeyType {
	return s.kty
}

func (s SignRequest) KeyType() service.KeyType {
	return s.keyType
}

func (s SignRequest) Data() []byte {
	return s.data
}

func (s SignRequest) UseEnsResolution() bool {
	return s.useEnsResolution
}
