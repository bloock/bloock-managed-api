package request

import "github.com/bloock/bloock-sdk-go/v2/entity/key"

type CreateManagedKeyRequest struct {
	name       string
	keyType    key.KeyType
	expiration int
	level      key.KeyProtectionLevel
}

func NewCreateManagedKeyRequest(name string, keyType key.KeyType, expiration int, level key.KeyProtectionLevel) *CreateManagedKeyRequest {
	return &CreateManagedKeyRequest{name: name, keyType: keyType, expiration: expiration, level: level}
}

func (c CreateManagedKeyRequest) Name() string {
	return c.name
}

func (c CreateManagedKeyRequest) KeyType() key.KeyType {
	return c.keyType
}

func (c CreateManagedKeyRequest) Expiration() int {
	return c.expiration
}

func (c CreateManagedKeyRequest) Level() key.KeyProtectionLevel {
	return c.level
}
