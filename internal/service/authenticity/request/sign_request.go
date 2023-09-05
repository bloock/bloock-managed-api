package request

import (
	"bloock-managed-api/internal/domain"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type SignRequest struct {
	publicKey        string
	privateKey       *string
	keySource        domain.KeyType
	keyID            uuid.UUID
	keyType          key.KeyType
	data             []byte
	useEnsResolution bool
}

func NewSignRequest(publicKey string, privateKey *string, keySource domain.KeyType, keyID uuid.UUID, keyType key.KeyType, data []byte, useEnsResolution bool) *SignRequest {
	return &SignRequest{publicKey: publicKey, privateKey: privateKey, keySource: keySource, keyID: keyID, keyType: keyType, data: data, useEnsResolution: useEnsResolution}
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

func (s SignRequest) KeyType() key.KeyType {
	return s.keyType
}

func (s SignRequest) KeySource() domain.KeyType {
	return s.keySource
}

func (s SignRequest) Data() []byte {
	return s.data
}

func (s SignRequest) UseEnsResolution() bool {
	return s.useEnsResolution
}
