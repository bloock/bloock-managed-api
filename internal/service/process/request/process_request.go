package request

import (
	"bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"strconv"
)

type ProcessRequest struct {
	file                  []byte
	isIntegrityEnabled    bool
	isAuthenticityEnabled bool
	keyType               domain.KeyType
	kty                   key.KeyType
	keyID                 uuid.UUID
	hostingType           domain.HostingType
	useEnsResolution      bool
}

func NewProcessRequest(data []byte, integrityEnabled string, authenticityEnabled string, keyType string, kty string, kid string, availabilityType string, ensRes string) (*ProcessRequest, error) {
	processRequestInstance := &ProcessRequest{}

	isIntegrityEnabled, err := strconv.ParseBool(integrityEnabled)
	if err != nil {
		return nil, err
	}

	isAuthenticityEnabled, err := strconv.ParseBool(authenticityEnabled)
	if err != nil {
		return nil, err
	}
	processRequestInstance.isAuthenticityEnabled = isAuthenticityEnabled
	if isAuthenticityEnabled {

		keyID, err := uuid.Parse(kid)
		if err != nil {
			return nil, err
		}
		processRequestInstance.keyID = keyID
		ktyp, err := domain.ValidateKeyType(kty)
		if err != nil {
			return nil, err
		}
		processRequestInstance.kty = ktyp
		authenticityKeyType, err := domain.ParseKeyType(keyType)
		if err != nil {
			return nil, err
		}
		processRequestInstance.keyType = authenticityKeyType

		useEnsResolution, err := strconv.ParseBool(ensRes)
		if err != nil {
			return nil, err
		}
		processRequestInstance.useEnsResolution = useEnsResolution

	}

	hostingType, err := domain.ParseHostingType(availabilityType)
	if err != nil {
		return nil, err
	}

	processRequestInstance.file = data
	processRequestInstance.hostingType = hostingType
	processRequestInstance.isIntegrityEnabled = isIntegrityEnabled
	return &ProcessRequest{
		file:                  data,
		isIntegrityEnabled:    isIntegrityEnabled,
		isAuthenticityEnabled: isAuthenticityEnabled,
		hostingType:           hostingType,
	}, nil
}

func (s ProcessRequest) KeyID() uuid.UUID {
	return s.keyID
}

func (s ProcessRequest) UseEnsResolution() bool {
	return s.useEnsResolution
}

func (s ProcessRequest) Data() []byte {
	return s.file
}
func (s ProcessRequest) ReplaceDataWith(newData []byte) {
	s.file = newData
}
func (s ProcessRequest) IsAuthenticityEnabled() bool {
	return s.isAuthenticityEnabled
}

func (s ProcessRequest) IsIntegrityEnabled() bool {
	return s.isIntegrityEnabled
}

func (s ProcessRequest) HostingType() domain.HostingType {
	return s.hostingType
}

func (s ProcessRequest) KeyType() domain.KeyType {
	return s.keyType
}

func (s ProcessRequest) Kty() key.KeyType {
	return s.kty
}
