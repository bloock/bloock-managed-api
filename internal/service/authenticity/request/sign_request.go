package request

import (
	"bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type SignRequest struct {
	publicKey        string
	privateKey       *string
	keyID            uuid.UUID
	kty              key.KeyType
	keyType          domain.KeyType
	data             []byte
	useEnsResolution bool
}

func NewSignRequest(publicKey string, privateKey *string, keyID uuid.UUID, kty key.KeyType, keyType domain.KeyType, data []byte, useEnsResolution bool) *SignRequest {
	return &SignRequest{publicKey: publicKey, privateKey: privateKey, keyID: keyID, kty: kty, keyType: keyType, data: data, useEnsResolution: useEnsResolution}
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

func (s SignRequest) KeyType() domain.KeyType {
	return s.keyType
}

func (s SignRequest) Data() []byte {
	return s.data
}

func (s SignRequest) UseEnsResolution() bool {
	return s.useEnsResolution
}
