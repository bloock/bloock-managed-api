package request

import (
	"bloock-managed-api/internal/domain"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type EncryptRequest struct {
	publicKey  string
	privateKey *string
	keySource  domain.KeyType
	keyID      uuid.UUID
	keyType    key.KeyType
	data       []byte
}

func NewEncryptRequest(publicKey string, privateKey *string, keySource domain.KeyType, keyID uuid.UUID, keyType key.KeyType, data []byte) EncryptRequest {
	return EncryptRequest{publicKey: publicKey, privateKey: privateKey, keySource: keySource, keyID: keyID, keyType: keyType, data: data}
}

func (s EncryptRequest) PublicKey() string {
	return s.publicKey
}

func (s EncryptRequest) PrivateKey() *string {
	return s.privateKey
}

func (s EncryptRequest) KeyID() uuid.UUID {
	return s.keyID
}

func (s EncryptRequest) KeyType() key.KeyType {
	return s.keyType
}

func (s EncryptRequest) KeySource() domain.KeyType {
	return s.keySource
}

func (s EncryptRequest) Data() []byte {
	return s.data
}
