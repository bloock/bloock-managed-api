package request

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"errors"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type ProcessRequest struct {
	file                  []byte
	isIntegrityEnabled    bool
	isAuthenticityEnabled bool
	keySource             domain.KeyType
	keyID                 uuid.UUID
	keyType               key.KeyType
	hostingType           domain.HostingType
	useEnsResolution      bool
}

func NewProcessRequest(data []byte, integrityEnabled bool, authenticityEnabled bool, keySource string, keyType string, kid string, useEns bool, availabilityType string) (*ProcessRequest, error) {
	processRequestInstance := &ProcessRequest{}

	processRequestInstance.file = data
	processRequestInstance.isIntegrityEnabled = integrityEnabled
	processRequestInstance.isAuthenticityEnabled = authenticityEnabled

	if authenticityEnabled {
		authenticityKeySource, err := domain.ParseKeySource(keySource)
		if err != nil {
			return nil, err
		}
		processRequestInstance.keySource = authenticityKeySource

		kty, err := domain.ValidateKeyType(keyType)
		if err != nil {
			return nil, err
		}
		processRequestInstance.keyType = kty

		if authenticityKeySource == domain.MANAGED_KEY || authenticityKeySource == domain.MANAGED_CERTIFICATE {
			// Managed key or certificate

			keyID, err := uuid.Parse(kid)
			if err != nil {
				return nil, err
			}
			processRequestInstance.keyID = keyID
		} else {
			if config.Configuration.PublicKey == "" {
				return nil, errors.New("no public key loaded")
			}

			if config.Configuration.PrivateKey == "" {
				return nil, errors.New("no private key loaded")
			}
		}

		processRequestInstance.useEnsResolution = useEns
	}

	hostingType, err := domain.ParseHostingType(availabilityType)
	if err != nil {
		return nil, err
	}
	processRequestInstance.hostingType = hostingType

	return processRequestInstance, nil
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
func (s *ProcessRequest) ReplaceDataWith(newData []byte) {
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

func (s ProcessRequest) KeySource() domain.KeyType {
	return s.keySource
}

func (s ProcessRequest) KeyType() key.KeyType {
	return s.keyType
}
