package request

import "github.com/google/uuid"

type SignRequest struct {
	localKey   uuid.UUID
	commonName *string
	data       []byte
}

func NewSignRequest(localKey uuid.UUID, commonName *string, data []byte) *SignRequest {
	return &SignRequest{localKey: localKey, commonName: commonName, data: data}
}

func (s SignRequest) LocalKey() uuid.UUID {
	return s.localKey
}

func (s SignRequest) CommonName() *string {
	return s.commonName
}

func (s SignRequest) Data() []byte {
	return s.data
}
