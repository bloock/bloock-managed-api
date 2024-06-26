package request

import (
	"os"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/pkg/request"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type IntegrityRequest struct {
	Enabled   bool
	Aggregate bool
}

type LocalKeyRequest struct {
	KeyType key.KeyType
	Key     string
}

type LocalCertificateRequest struct {
	Pkcs12        []byte
	Pkcs12Pasword string
}

type ManagedKeyRequest struct {
	Uuid uuid.UUID
}

type ManagedCertificateRequest struct {
	Uuid uuid.UUID
}

type AccessControlRequest struct {
	AccessControlType domain.AccessControlType
	AccessCode        string
}

type AuthenticityRequest struct {
	Enabled            bool
	KeySource          domain.KeyType
	LocalKey           *LocalKeyRequest
	LocalCertificate   *LocalCertificateRequest
	ManagedKey         *ManagedKeyRequest
	ManagedCertificate *ManagedCertificateRequest
	AccessControl      *AccessControlRequest
}

type EncryptionRequest struct {
	Enabled            bool
	KeySource          domain.KeyType
	LocalKey           *LocalKeyRequest
	LocalCertificate   *LocalCertificateRequest
	ManagedKey         *ManagedKeyRequest
	ManagedCertificate *ManagedCertificateRequest
	AccessControl      *AccessControlRequest
}

type AvailabilityRequest struct {
	Enabled     bool
	Hostingtype domain.HostingType
}

type ProcessRequest struct {
	File         domain.File
	Integrity    IntegrityRequest
	Authenticity AuthenticityRequest
	Encryption   EncryptionRequest
	Availability AvailabilityRequest
}

func NewProcessRequest(file domain.File, request request.ProcessFormRequest) (*ProcessRequest, error) {
	processRequestInstance := &ProcessRequest{}

	processRequestInstance.File = file

	processRequestInstance.Integrity.Aggregate = false
	if request.Integrity.Enabled {
		integrityRequest := IntegrityRequest{
			Enabled:   request.Integrity.Enabled,
			Aggregate: request.Integrity.Aggregate,
		}
		processRequestInstance.Integrity = integrityRequest
	}

	if request.Authenticity.Enabled {
		authenticityRequest := AuthenticityRequest{
			Enabled: request.Authenticity.Enabled,
		}

		authenticityKeySource, err := domain.ParseKeySource(request.Authenticity.KeySource)
		if err != nil {
			return nil, err
		}
		authenticityRequest.KeySource = authenticityKeySource

		switch authenticityKeySource {
		case domain.LOCAL_KEY:
			kty, err := domain.ValidateKeyType(config.Configuration.Authenticity.KeyConfig.KeyType)
			if err != nil {
				return nil, err
			}

			authenticityRequest.LocalKey = &LocalKeyRequest{
				KeyType: kty,
				Key:     config.Configuration.Authenticity.KeyConfig.Key,
			}
		case domain.MANAGED_KEY:
			keyID, err := uuid.Parse(request.Authenticity.Key)
			if err != nil {
				return nil, err
			}
			authenticityRequest.ManagedKey = &ManagedKeyRequest{
				Uuid: keyID,
			}
		case domain.LOCAL_CERTIFICATE:
			pkcs12, err := os.ReadFile(config.Configuration.Authenticity.CertificateConfig.Pkcs12Path)
			if err != nil {
				return nil, err
			}
			authenticityRequest.LocalCertificate = &LocalCertificateRequest{
				Pkcs12:        pkcs12,
				Pkcs12Pasword: config.Configuration.Authenticity.CertificateConfig.Pkcs12Password,
			}
		case domain.MANAGED_CERTIFICATE:
			keyID, err := uuid.Parse(request.Authenticity.Key)
			if err != nil {
				return nil, err
			}
			authenticityRequest.ManagedCertificate = &ManagedCertificateRequest{
				Uuid: keyID,
			}
		}

		if request.Authenticity.AccessEnabled {
			if request.Authenticity.AccessCode == "" {
				return nil, domain.ErrEmptyAccessCode
			}
			authenticityAccessType, err := domain.ParseAccessControlType(request.Authenticity.AccessType)
			if err != nil {
				return nil, err
			}
			accessControl := &AccessControlRequest{
				AccessControlType: authenticityAccessType,
			}
			accessControl.AccessCode = request.Authenticity.AccessCode
			authenticityRequest.AccessControl = accessControl
		} else {
			authenticityRequest.AccessControl = nil
		}

		processRequestInstance.Authenticity = authenticityRequest
	}

	if request.Encryption.Enabled {
		encryptionRequest := EncryptionRequest{
			Enabled: request.Encryption.Enabled,
		}

		encryptionKeySource, err := domain.ParseKeySource(request.Encryption.KeySource)
		if err != nil {
			return nil, err
		}
		encryptionRequest.KeySource = encryptionKeySource

		switch encryptionKeySource {
		case domain.LOCAL_KEY:
			kty, err := domain.ValidateKeyType(config.Configuration.Encryption.KeyConfig.KeyType)
			if err != nil {
				return nil, err
			}

			encryptionRequest.LocalKey = &LocalKeyRequest{
				KeyType: kty,
				Key:     config.Configuration.Encryption.KeyConfig.Key,
			}
		case domain.MANAGED_KEY:
			keyID, err := uuid.Parse(request.Encryption.Key)
			if err != nil {
				return nil, err
			}
			encryptionRequest.ManagedKey = &ManagedKeyRequest{
				Uuid: keyID,
			}
		case domain.LOCAL_CERTIFICATE:
			pkcs12, err := os.ReadFile(config.Configuration.Encryption.CertificateConfig.Pkcs12Path)
			if err != nil {
				return nil, err
			}
			encryptionRequest.LocalCertificate = &LocalCertificateRequest{
				Pkcs12:        pkcs12,
				Pkcs12Pasword: config.Configuration.Encryption.CertificateConfig.Pkcs12Password,
			}
		case domain.MANAGED_CERTIFICATE:
			keyID, err := uuid.Parse(request.Encryption.Key)
			if err != nil {
				return nil, err
			}
			encryptionRequest.ManagedCertificate = &ManagedCertificateRequest{
				Uuid: keyID,
			}
		}

		if request.Encryption.AccessEnabled {
			if request.Encryption.AccessCode == "" {
				return nil, domain.ErrEmptyAccessCode
			}
			encryptionAccessType, err := domain.ParseAccessControlType(request.Encryption.AccessType)
			if err != nil {
				return nil, err
			}
			accessControl := &AccessControlRequest{
				AccessControlType: encryptionAccessType,
			}
			accessControl.AccessCode = request.Encryption.AccessCode
			encryptionRequest.AccessControl = accessControl
		} else {
			encryptionRequest.AccessControl = nil
		}

		processRequestInstance.Encryption = encryptionRequest
	}

	if request.Availability.Enabled {
		hostingType, err := domain.ParseHostingType(request.Availability.Type)
		if err != nil {
			return nil, err
		}
		processRequestInstance.Availability = AvailabilityRequest{
			Enabled:     request.Availability.Enabled,
			Hostingtype: hostingType,
		}
	}

	return processRequestInstance, nil
}
