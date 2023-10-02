package request

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"errors"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type ProcessRequest struct {
	file                         []byte
	integrityEnabled             bool
	authenticityEnabled          bool
	authenticityKeySource        domain.KeyType
	authenticityKeyID            uuid.UUID
	authenticityKeyType          key.KeyType
	authenticityUseEnsResolution bool
	encryptionEnabled            bool
	encryptionKeySource          domain.KeyType
	encryptionKeyID              uuid.UUID
	encryptionKeyType            key.KeyType
	hostingType                  domain.HostingType
}

func NewProcessRequest(data []byte, integrityEnabled bool, authenticityEnabled bool, authenticityKeySource string, authenticityKeyType string, authenticityKid string, authenticityUseEns bool, encryptionEnabled bool, encryptionKeySource string, encryptionKeyType string, encryptionKid string, availabilityType string) (*ProcessRequest, error) {
	processRequestInstance := &ProcessRequest{}

	processRequestInstance.file = data
	processRequestInstance.integrityEnabled = integrityEnabled

	processRequestInstance.authenticityEnabled = authenticityEnabled
	if authenticityEnabled {
		authenticityKeySource, err := domain.ParseKeySource(authenticityKeySource)
		if err != nil {
			return nil, err
		}
		processRequestInstance.authenticityKeySource = authenticityKeySource

		kty, err := domain.ValidateKeyType(authenticityKeyType)
		if err != nil {
			return nil, err
		}
		processRequestInstance.authenticityKeyType = kty

		if authenticityKeySource == domain.MANAGED_KEY || authenticityKeySource == domain.MANAGED_CERTIFICATE {
			// Managed key or certificate

			keyID, err := uuid.Parse(authenticityKid)
			if err != nil {
				return nil, err
			}
			processRequestInstance.authenticityKeyID = keyID
		} else {
			if config.Configuration.AuthenticityPublicKey == "" {
				return nil, errors.New("no public key loaded")
			}

			if config.Configuration.AuthenticityPrivateKey == "" {
				return nil, errors.New("no private key loaded")
			}
		}

		processRequestInstance.authenticityUseEnsResolution = authenticityUseEns
	}

	processRequestInstance.encryptionEnabled = encryptionEnabled
	if encryptionEnabled {
		encryptionKeySource, err := domain.ParseKeySource(encryptionKeySource)
		if err != nil {
			return nil, err
		}
		processRequestInstance.encryptionKeySource = encryptionKeySource

		kty, err := domain.ValidateKeyType(encryptionKeyType)
		if err != nil {
			return nil, err
		}
		processRequestInstance.encryptionKeyType = kty

		if encryptionKeySource == domain.MANAGED_KEY || encryptionKeySource == domain.MANAGED_CERTIFICATE {
			// Managed key or certificate

			encryptionKeyID, err := uuid.Parse(encryptionKid)
			if err != nil {
				return nil, err
			}
			processRequestInstance.encryptionKeyID = encryptionKeyID
		} else {
			if config.Configuration.EncryptionPublicKey == "" {
				return nil, errors.New("no public key loaded")
			}
		}

		processRequestInstance.authenticityUseEnsResolution = authenticityUseEns
	}

	hostingType, err := domain.ParseHostingType(availabilityType)
	if err != nil {
		return nil, err
	}
	processRequestInstance.hostingType = hostingType

	return processRequestInstance, nil
}

func (s ProcessRequest) Data() []byte {
	return s.file
}
func (s *ProcessRequest) ReplaceDataWith(newData []byte) {
	s.file = newData
}

func (s ProcessRequest) IsIntegrityEnabled() bool {
	return s.integrityEnabled
}

func (s ProcessRequest) IsAuthenticityEnabled() bool {
	return s.authenticityEnabled
}

func (s ProcessRequest) AuthenticityKeySource() domain.KeyType {
	return s.authenticityKeySource
}

func (s ProcessRequest) AuthenticityKeyType() key.KeyType {
	return s.authenticityKeyType
}

func (s ProcessRequest) AuthenticityKeyID() uuid.UUID {
	return s.authenticityKeyID
}

func (s ProcessRequest) AuthenticityUseEnsResolution() bool {
	return s.authenticityUseEnsResolution
}

func (s ProcessRequest) IsEncryptionEnabled() bool {
	return s.encryptionEnabled
}

func (s ProcessRequest) EncryptionKeySource() domain.KeyType {
	return s.encryptionKeySource
}

func (s ProcessRequest) EncryptionKeyType() key.KeyType {
	return s.encryptionKeyType
}

func (s ProcessRequest) EncryptionKeyID() uuid.UUID {
	return s.encryptionKeyID
}

func (s ProcessRequest) HostingType() domain.HostingType {
	return s.hostingType
}
