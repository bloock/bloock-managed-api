package process

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/service/process/request"
	"github.com/bloock/bloock-managed-api/internal/service/process/response"
	"github.com/rs/zerolog"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

var (
	ErrSignKeyNotSupported    = errors.New("key type not supported for signing")
	ErrEncryptKeyNotSupported = errors.New("key type not supported for encrypting")
	ErrUnsupportedHosting     = errors.New("unsupported hosting type")
)

type ProcessService struct {
	integrityRepository    repository.IntegrityRepository
	keyRepository          repository.KeyRepository
	authenticityRepository repository.AuthenticityRepository
	encryptionRepository   repository.EncryptionRepository
	availabilityRepository repository.AvailabilityRepository
	metadataRepository     repository.MetadataRepository
	notificationRepository repository.NotificationRepository
}

func NewProcessService(ctx context.Context, logger zerolog.Logger) *ProcessService {
	return &ProcessService{
		integrityRepository:    bloock_repository.NewBloockIntegrityRepository(ctx, logger),
		keyRepository:          bloock_repository.NewBloockKeyRepository(ctx, logger),
		authenticityRepository: bloock_repository.NewBloockAuthenticityRepository(ctx, logger),
		encryptionRepository:   bloock_repository.NewBloockEncryptionRepository(ctx, logger),
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, logger),
		metadataRepository:     bloock_repository.NewBloockMetadataRepository(ctx, logger),
		notificationRepository: bloock_repository.NewHttpNotificationRepository(ctx, logger),
	}
}

func (s ProcessService) LoadUrl(ctx context.Context, url *url.URL) (domain.File, error) {
	fileBytes, err := s.availabilityRepository.FindFile(ctx, url.String())
	if err != nil {
		return domain.File{}, err
	}

	filename := path.Base(url.Path)
	if filename == "" {
		pathParts := strings.Split(url.Path, "/")

		// If it's empty, use the second-to-last part as the filename
		if len(pathParts) >= 2 {
			filename = pathParts[len(pathParts)-2]
		}
	}

	return domain.NewFile(fileBytes, filename, http.DetectContentType(fileBytes)), nil
}

func (s ProcessService) Process(ctx context.Context, req request.ProcessRequest) (*response.ProcessResponse, error) {
	record, err := s.metadataRepository.GetRecord(ctx, req.File.Bytes())
	if err != nil {
		return nil, err
	}

	fileHash, err := record.GetHash()
	if err != nil {
		return nil, err
	}

	certification := domain.Certification{
		Data:   req.File.Bytes(),
		Hash:   fileHash,
		Record: record,
	}

	responseBuilder := response.NewProcessResponseBuilder()

	if req.Authenticity.Enabled {
		_, _, record, err := s.sign(ctx, &req.File, &req.Authenticity)
		if err != nil {
			return nil, err
		}

		newHash, err := record.GetHash()
		if err != nil {
			return nil, err
		}

		req.File.File = record.Retrieve()

		certification.Data = record.Retrieve()
		certification.Record = record
		certification.Hash = newHash

		rd, err := s.metadataRepository.GetRecordDetails(ctx, certification.Data)
		signatures := make([]response.Signature, 0)
		if rd.AuthenticityDetails != nil {
			for _, sig := range rd.AuthenticityDetails.Signatures {
				signatures = append(signatures, response.NewSignature(sig.Signature, sig.Alg, sig.Kid, sig.MessageHash, sig.Subject))
			}
		}
		responseBuilder.SignResponse(*response.NewSignResponse(signatures))
	}

	if req.Integrity.Enabled {
		newCertification, err := s.certify(ctx, certification.Data)
		if err != nil {
			return nil, err
		}
		certification = newCertification

		responseBuilder.CertificationResponse(*response.NewIntegrityResponse(certification.Hash, certification.AnchorID))
	}

	if req.Encryption.Enabled {
		_, encryptedRecord, err := s.encrypt(ctx, &req.File, &req.Encryption)
		if err != nil {
			return nil, err
		}

		req.File.File = encryptedRecord.Retrieve()
		certification.Data = encryptedRecord.Retrieve()
		certification.Record = encryptedRecord

		rd, err := s.metadataRepository.GetRecordDetails(ctx, certification.Data)
		if rd.EncryptionDetails != nil {
			responseBuilder.EncryptResponse(*response.NewEncryptResponse(rd.EncryptionDetails.Key, rd.EncryptionDetails.Alg, rd.EncryptionDetails.Subject))
		}
	}

	if req.Availability.Enabled {
		dataID, err := s.upload(ctx, &req.File, *certification.Record, &req.Availability)
		if err != nil {
			return nil, err
		}
		certification.DataID = dataID

		rd, err := s.metadataRepository.GetRecordDetails(ctx, certification.Data)
		var contentType string
		if rd.AvailabilityDetails.ContentType != nil {
			contentType = *rd.AvailabilityDetails.ContentType
		}

		responseBuilder.AvailabilityResponse(*response.NewAvailabilityResponse(certification.DataID, req.Availability.Hostingtype, contentType, rd.AvailabilityDetails.Size))
	} else {
		if req.Integrity.Enabled {
			if _, err = s.availabilityRepository.UploadTmp(ctx, &req.File, *certification.Record); err != nil {
				return nil, err
			}
		}
	}

	if err = s.metadataRepository.UpdateCertification(ctx, certification); err != nil {
		return nil, err
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

func (s ProcessService) sign(ctx context.Context, file *domain.File, request *request.AuthenticityRequest) (*string, string, *record.Record, error) {
	switch request.KeySource {
	case domain.LOCAL_KEY:
		localKey, err := s.keyRepository.LoadLocalKey(ctx, request.LocalKey.KeyType, request.LocalKey.Key)
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithLocalKey(ctx, file.Bytes(), *localKey)
		if err != nil {
			return nil, "", nil, err
		}
		return &localKey.Key, signature, record, nil

	case domain.MANAGED_KEY:
		managedKey, err := s.keyRepository.LoadManagedKey(ctx, request.ManagedKey.Uuid.String())
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithManagedKey(ctx, file.Bytes(), *managedKey)
		if err != nil {
			return nil, "", nil, err
		}
		return &managedKey.ID, signature, record, nil

	case domain.LOCAL_CERTIFICATE:
		localCertificate, err := s.keyRepository.LoadLocalCertificate(ctx, request.LocalCertificate.Pkcs12, request.LocalCertificate.Pkcs12Pasword)
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithLocalCertificate(ctx, file.Bytes(), *localCertificate)
		if err != nil {
			return nil, "", nil, err
		}

		return &signature, signature, record, nil

	case domain.MANAGED_CERTIFICATE:
		managedCertificate, err := s.keyRepository.LoadManagedCertificate(ctx, request.ManagedCertificate.Uuid.String())
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithManagedCertificate(ctx, file.Bytes(), *managedCertificate)
		if err != nil {
			return nil, "", nil, err
		}
		return &managedCertificate.ID, signature, record, nil
	}

	return nil, "", nil, ErrSignKeyNotSupported
}

func (s ProcessService) encrypt(ctx context.Context, file *domain.File, request *request.EncryptionRequest) (*string, *record.Record, error) {
	switch request.KeySource {
	case domain.LOCAL_KEY:
		localKey, err := s.keyRepository.LoadLocalKey(ctx, request.LocalKey.KeyType, request.LocalKey.Key)
		if err != nil {
			return nil, nil, err
		}

		record, err := s.encryptionRepository.EncryptWithLocalKey(ctx, file.Bytes(), *localKey)
		if err != nil {
			return nil, record, err
		}
		return &localKey.Key, record, nil

	case domain.MANAGED_KEY:
		managedKey, err := s.keyRepository.LoadManagedKey(ctx, request.ManagedKey.Uuid.String())
		if err != nil {
			return nil, nil, err
		}

		record, err := s.encryptionRepository.EncryptWithManagedKey(ctx, file.Bytes(), *managedKey)
		if err != nil {
			return nil, record, err
		}
		return &managedKey.ID, record, nil

	case domain.LOCAL_CERTIFICATE:
		return nil, nil, ErrEncryptKeyNotSupported
	case domain.MANAGED_CERTIFICATE:
		managedCertificate, err := s.keyRepository.LoadManagedCertificate(ctx, request.ManagedCertificate.Uuid.String())
		if err != nil {
			return nil, nil, err
		}

		record, err := s.encryptionRepository.EncryptWithManagedCertificate(ctx, file.Bytes(), *managedCertificate)
		if err != nil {
			return nil, record, err
		}
		return &managedCertificate.ID, record, nil
	}

	return nil, nil, ErrEncryptKeyNotSupported
}

func (a ProcessService) upload(ctx context.Context, file *domain.File, record record.Record, request *request.AvailabilityRequest) (string, error) {
	switch request.Hostingtype {
	case domain.HOSTED:
		hostedID, err := a.availabilityRepository.UploadHosted(ctx, file, record)
		if err != nil {
			return "", err
		}
		return hostedID, err
	case domain.IPFS:
		ipfsID, err := a.availabilityRepository.UploadIpfs(ctx, file, record)
		if err != nil {
			return "", err
		}
		return ipfsID, err
	case domain.LOCAL:
		path, err := a.availabilityRepository.UploadLocal(ctx, file)
		if err != nil {
			return "", err
		}
		return path, err
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
