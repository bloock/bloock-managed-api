package process

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

var (
	ErrSignKeyNotSupported    = errors.New("key type not supported for signing")
	ErrEncryptKeyNotSupported = errors.New("key type not supported for encrypting")
	ErrUnsupportedHosting     = errors.New("unsupported hosting type")
)

type ProcessService struct {
	integrityRepository    repository.IntegrityRepository
	authenticityRepository repository.AuthenticityRepository
	encryptionRepository   repository.EncryptionRepository
	availabilityRepository repository.AvailabilityRepository
	metadataRepository     repository.MetadataRepository
	notificationRepository repository.NotificationRepository
}

func NewProcessService(
	integrityRepository repository.IntegrityRepository,
	authenticityRepository repository.AuthenticityRepository,
	encryptionRepository repository.EncryptionRepository,
	availabilityRepository repository.AvailabilityRepository,
	metadataRepository repository.MetadataRepository,
	notificationRepository repository.NotificationRepository,
) *ProcessService {

	return &ProcessService{
		integrityRepository:    integrityRepository,
		authenticityRepository: authenticityRepository,
		encryptionRepository:   encryptionRepository,
		availabilityRepository: availabilityRepository,
		metadataRepository:     metadataRepository,
		notificationRepository: notificationRepository,
	}
}

func (s ProcessService) Process(ctx context.Context, req request.ProcessRequest) (*response.ProcessResponse, error) {
	var file []byte
	if req.File() != nil {
		file = req.File()
	} else if req.URL() != "" {
		fileBytes, err := s.availabilityRepository.FindFile(ctx, req.URL())
		if err != nil {
			return nil, err
		}
		req = req.SetContentType(http.DetectContentType(fileBytes))
		file = fileBytes
	} else {
		return nil, errors.New("you must provide a file or URL")
	}

	record, err := s.metadataRepository.GetRecord(ctx, file)
	if err != nil {
		return nil, err
	}

	fileHash, err := record.GetHash()
	if err != nil {
		return nil, err
	}

	certification := domain.Certification{
		Data:   file,
		Hash:   fileHash,
		Record: record,
	}

	responseBuilder := response.NewProcessResponseBuilder()

	if req.IsAuthenticityEnabled() {
		signRequest := request.NewSignRequest(
			config.Configuration.AuthenticityPublicKey,
			&config.Configuration.AuthenticityPrivateKey,
			req.AuthenticityKeySource(),
			req.AuthenticityKeyID(),
			req.AuthenticityKeyType(),
			req.AuthenticityUseEnsResolution(),
			file,
		)
		signature, record, err := s.sign(ctx, signRequest)
		if err != nil {
			return nil, err
		}

		newHash, err := record.GetHash()
		if err != nil {
			return nil, err
		}

		certification.Data = record.Retrieve()
		certification.Record = record
		certification.Hash = newHash
		responseBuilder.SignResponse(*response.NewSignResponse(signature))
	}

	if req.IsIntegrityEnabled() {
		newCertification, err := s.certify(ctx, certification.Data)
		if err != nil {
			return nil, err
		}
		certification = newCertification
		responseBuilder.CertificationResponse(*response.NewIntegrityResponse(certification.Hash, certification.AnchorID))
	}

	if req.IsEncryptionEnabled() {
		encryptRequest := request.NewEncryptRequest(
			config.Configuration.EncryptionPublicKey,
			&config.Configuration.EncryptionPrivateKey,
			req.EncryptionKeySource(),
			req.EncryptionKeyID(),
			req.EncryptionKeyType(),
			file,
		)
		encryptedRecord, err := s.encrypt(ctx, encryptRequest)
		if err != nil {
			return nil, err
		}

		newHash, err := record.GetHash()
		if err != nil {
			return nil, err
		}

		certification.Data = encryptedRecord.Retrieve()
		certification.Record = encryptedRecord
		certification.Hash = newHash
	}

	if req.HostingType() != domain.NONE {
		dataID, err := s.upload(ctx, fmt.Sprintf("%s%s", req.Filename(), req.FileExtension()), certification.Record, req.HostingType())
		if err != nil {
			return nil, err
		}
		certification.DataID = dataID
		if err = s.metadataRepository.UpdateCertification(ctx, certification); err != nil {
			return nil, err
		}
		responseBuilder.AvailabilityResponse(*response.NewAvailabilityResponse(certification.DataID, req.HostingType()))
	} else {
		if req.IsIntegrityEnabled() {
			if _, err = s.availabilityRepository.UploadTmp(ctx, certification.Record); err != nil {
				return nil, err
			}
		}
	}
	responseBuilder.HashResponse(certification.Hash)

	if certification.AnchorID == 0 {
		if err = s.notify(ctx, []domain.Certification{certification}); err != nil {
			return nil, err
		}
	}

	return responseBuilder.Build(), nil
}

func (s ProcessService) certify(ctx context.Context, data []byte) (domain.Certification, error) {
	certification, err := s.integrityRepository.Certify(ctx, data)
	if err != nil {
		return domain.Certification{}, err
	}

	if err = s.metadataRepository.SaveCertification(ctx, certification); err != nil {
		return domain.Certification{}, err
	}

	return certification, nil
}

func (s ProcessService) sign(ctx context.Context, request request.SignRequest) (string, *record.Record, error) {
	switch request.KeySource() {
	case domain.LOCAL_KEY:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithLocalKeyEns(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}

	case domain.MANAGED_KEY:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKeyEns(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}

	case domain.LOCAL_CERTIFICATE:
		return "", nil, nil
	case domain.MANAGED_CERTIFICATE:
		if request.KeyType() == key.EcP256k && !request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}

		if request.KeyType() == key.EcP256k && request.UseEnsResolution() {
			signature, record, err := s.authenticityRepository.
				SignECWithManagedKeyEns(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return "", nil, err
			}
			return signature, record, nil
		}
	}

	return "", nil, ErrSignKeyNotSupported
}

func (s ProcessService) encrypt(ctx context.Context, request request.EncryptRequest) (*record.Record, error) {
	switch request.KeySource() {
	case domain.LOCAL_KEY:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey(), request.PrivateKey())
			if err != nil {
				return record, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithLocalKey(ctx, request.Data(), request.KeyType(), request.PublicKey())
			if err != nil {
				return nil, err
			}

			return record, nil
		}

	case domain.MANAGED_KEY:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		}

	case domain.LOCAL_CERTIFICATE:
		return nil, nil
	case domain.MANAGED_CERTIFICATE:
		switch request.KeyType() {
		case key.Rsa2048, key.Rsa3072, key.Rsa4096:
			record, err := s.encryptionRepository.EncryptRSAWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		case key.Aes128, key.Aes256:
			record, err := s.encryptionRepository.EncryptAESWithManagedKey(ctx, request.Data(), request.KeyID().String())
			if err != nil {
				return nil, err
			}

			return record, nil
		}
	}

	return nil, ErrEncryptKeyNotSupported
}

func (a ProcessService) upload(ctx context.Context, filename string, record *record.Record, hostingType domain.HostingType) (string, error) {
	switch hostingType {
	case domain.HOSTED:
		hostedID, err := a.availabilityRepository.UploadHosted(ctx, record)
		if err != nil {
			return "", err
		}
		return hostedID, err
	case domain.IPFS:
		ipfsID, err := a.availabilityRepository.UploadIpfs(ctx, record)
		if err != nil {
			return "", err
		}
		return ipfsID, err
	case domain.LOCAL:
		path, err := a.availabilityRepository.UploadLocal(ctx, filename, record)
		if err != nil {
			return "", err
		}
		return path, err
	case domain.NONE:
		return "", nil
	default:
		return "", ErrUnsupportedHosting
	}
}

func (n ProcessService) notify(ctx context.Context, certifications []domain.Certification) error {
	for _, crt := range certifications {
		var fileBytes []byte
		var err error

		if len(crt.Data) != 0 {
			fileBytes = crt.Data
		} else {
			if crt.DataID != "" {
				fileBytes, err = n.availabilityRepository.FindFile(ctx, crt.DataID)
				if err != nil {
					return err
				}
			} else {
				fileBytes, err = n.availabilityRepository.RetrieveTmp(ctx, crt.Hash)
				if err != nil {
					return err
				}
			}
		}

		if err = n.notificationRepository.NotifyCertification(crt.Hash, fileBytes); err != nil {
			return err
		}
	}

	return nil
}
