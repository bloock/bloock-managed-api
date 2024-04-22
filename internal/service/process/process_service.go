package process

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/service/process/request"
	"github.com/bloock/bloock-managed-api/internal/service/process/response"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

var (
	ErrSignKeyNotSupported    = errors.New("key type not supported for signing")
	ErrEncryptKeyNotSupported = errors.New("key type not supported for encrypting")
	ErrUnsupportedHosting     = errors.New("unsupported hosting type")
	ErrAggregateModeDisabled  = errors.New("aggregate mode disabled")
)

type ProcessService struct {
	integrityRepository    repository.IntegrityRepository
	keyRepository          repository.KeyRepository
	authenticityRepository repository.AuthenticityRepository
	encryptionRepository   repository.EncryptionRepository
	availabilityRepository repository.AvailabilityRepository
	metadataRepository     repository.MetadataRepository
	notificationRepository repository.NotificationRepository
	messageAggregator      repository.MessageAggregatorRepository
	processRepository      repository.ProcessRepository
	logger                 zerolog.Logger
}

func NewProcessService(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *ProcessService {
	logger := l.With().Caller().Str("component", "process-service").Logger()

	return &ProcessService{
		integrityRepository:    bloock_repository.NewBloockIntegrityRepository(ctx, l),
		keyRepository:          bloock_repository.NewBloockKeyRepository(ctx, l),
		authenticityRepository: bloock_repository.NewBloockAuthenticityRepository(ctx, l),
		encryptionRepository:   bloock_repository.NewBloockEncryptionRepository(ctx, l),
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, l),
		metadataRepository:     bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		notificationRepository: bloock_repository.NewHttpNotificationRepository(ctx, l),
		messageAggregator:      bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		processRepository:      bloock_repository.NewProcessRepository(ctx, l, ent),
		logger:                 logger,
	}
}

func (s ProcessService) Process(ctx context.Context, req request.ProcessRequest) (*response.ProcessResponse, error) {
	rec, err := s.metadataRepository.GetRecord(ctx, req.File.Bytes())
	if err != nil {
		log.Println("Enter")
		return nil, err
	}

	fileHash, err := rec.GetHash()
	if err != nil {
		return nil, err
	}

	certification := domain.Certification{
		Data:   req.File.Bytes(),
		Hash:   fileHash,
		Record: rec,
	}

	responseBuilder := response.NewProcessResponseBuilder()

	if req.Authenticity.Enabled {
		_, _, rec, err := s.sign(ctx, &req.File, &req.Authenticity)
		if err != nil {
			return nil, err
		}

		newHash, err := rec.GetHash()
		if err != nil {
			return nil, err
		}

		req.File.File = rec.Retrieve()

		certification.Data = rec.Retrieve()
		certification.Record = rec
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
		newCertification, err := s.certify(ctx, certification.Data, req.Integrity.Aggregate)
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

	if !req.Integrity.Aggregate {
		if err = s.metadataRepository.UpdateCertification(ctx, certification); err != nil {
			return nil, err
		}
	}

	responseBuilder.HashResponse(certification.Hash)
	responseBuilder.ProcessIDResponse(uuid.New().String())

	if certification.AnchorID == 0 && !req.Integrity.Aggregate {
		if err = s.notify(ctx, []domain.Certification{certification}); err != nil {
			return nil, err
		}
	}

	processResponse := responseBuilder.Build()

	if err = s.saveProcess(ctx, processResponse, req.File.FilenameWithExtension(), req.Integrity.Aggregate); err != nil {
		return nil, err
	}

	return processResponse, nil
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

func (s ProcessService) certify(ctx context.Context, data []byte, aggregate bool) (domain.Certification, error) {
	if aggregate {
		if !config.Configuration.Integrity.AggregateMode {
			return domain.Certification{}, ErrAggregateModeDisabled
		}
		rec, err := s.metadataRepository.GetRecord(ctx, data)
		if err != nil {
			return domain.Certification{}, err
		}
		hash, err := rec.GetHash()
		if err != nil {
			return domain.Certification{}, err
		}
		message := domain.Message{Hash: hash}
		if err = s.messageAggregator.SaveMessage(ctx, message); err != nil {
			return domain.Certification{}, err
		}
		return domain.Certification{Hash: hash, Data: data, Record: rec}, nil
	} else {
		certification, err := s.integrityRepository.Certify(ctx, data)
		if err != nil {
			return domain.Certification{}, err
		}
		if err = s.metadataRepository.SaveCertification(ctx, certification); err != nil {
			return domain.Certification{}, err
		}
		return certification, nil
	}
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

		accessControl, err := s.buildAccessControl(request.AccessControl)
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithManagedKey(ctx, file.Bytes(), *managedKey, accessControl)
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

		accessControl, err := s.buildAccessControl(request.AccessControl)
		if err != nil {
			return nil, "", nil, err
		}

		signature, record, err := s.authenticityRepository.
			SignWithManagedCertificate(ctx, file.Bytes(), *managedCertificate, accessControl)
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

		accessControl, err := s.buildAccessControl(request.AccessControl)
		if err != nil {
			return nil, nil, err
		}

		record, err := s.encryptionRepository.EncryptWithManagedKey(ctx, file.Bytes(), *managedKey, accessControl)
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

		accessControl, err := s.buildAccessControl(request.AccessControl)
		if err != nil {
			return nil, nil, err
		}

		record, err := s.encryptionRepository.EncryptWithManagedCertificate(ctx, file.Bytes(), *managedCertificate, accessControl)
		if err != nil {
			return nil, record, err
		}
		return &managedCertificate.ID, record, nil
	}

	return nil, nil, ErrEncryptKeyNotSupported
}

func (s ProcessService) upload(ctx context.Context, file *domain.File, record record.Record, request *request.AvailabilityRequest) (string, error) {
	switch request.Hostingtype {
	case domain.HOSTED:
		hostedID, err := s.availabilityRepository.UploadHosted(ctx, file, record)
		if err != nil {
			return "", err
		}
		return hostedID, err
	case domain.IPFS:
		ipfsID, err := s.availabilityRepository.UploadIpfs(ctx, file, record)
		if err != nil {
			return "", err
		}
		return ipfsID, err
	case domain.LOCAL:
		path, err := s.availabilityRepository.UploadLocal(ctx, file)
		if err != nil {
			return "", err
		}
		return path, err
	default:
		return "", ErrUnsupportedHosting
	}
}

func (s ProcessService) notify(ctx context.Context, certifications []domain.Certification) error {
	for _, crt := range certifications {
		var fileBytes []byte
		var err error

		if len(crt.Data) != 0 {
			fileBytes = crt.Data
		} else {
			if crt.DataID != "" {
				fileBytes, err = s.availabilityRepository.FindFile(ctx, crt.DataID)
				if err != nil {
					return err
				}
			} else {
				fileBytes, err = s.availabilityRepository.RetrieveTmp(ctx, crt.Hash)
				if err != nil {
					return err
				}
			}
		}

		if err = s.notificationRepository.NotifyCertification(crt.Hash, fileBytes); err != nil {
			return err
		}
	}

	return nil
}

func (s ProcessService) buildAccessControl(request *request.AccessControlRequest) (*key.AccessControl, error) {
	var accessControl key.AccessControl

	if request != nil {
		code := request.AccessCode
		switch request.AccessControlType {
		case domain.TotpAccessControl:
			accessControl.AccessControlTotp = key.NewAccessControlTotp(code)
		case domain.SecretAccessControl:
			accessControl.AccessControlSecret = key.NewAccessControlSecret(code)
		default:
			return nil, errors.New("invalid access control type")
		}
	} else {
		return nil, nil
	}

	return &accessControl, nil
}

func (s ProcessService) saveProcess(ctx context.Context, response *response.ProcessResponse, filename string, isAggregated bool) error {
	handlerResponse := response.MapToHandlerProcessResponse()

	newProcess := domain.Process{
		ID:              handlerResponse.ProcessID,
		Filename:        filename,
		ProcessResponse: handlerResponse,
	}

	return s.processRepository.SaveProcess(ctx, newProcess, isAggregated)
}
