package process

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	authenticity_request "bloock-managed-api/internal/service/authenticity/request"
	authenticity_response "bloock-managed-api/internal/service/authenticity/response"
	availability_response "bloock-managed-api/internal/service/availability/response"
	encryption_request "bloock-managed-api/internal/service/encryption/request"
	integrity_response "bloock-managed-api/internal/service/integrity/response"
	process_request "bloock-managed-api/internal/service/process/request"
	process_response "bloock-managed-api/internal/service/process/response"
	"context"
)

type ProcessService struct {
	integrityService    service.IntegrityService
	authenticityService service.AuthenticityService
	encryptionService   service.EncryptionService
	availabilityService service.AvailabilityService
	fileService         service.FileService
	notifyService       service.NotifyService
}

func NewProcessService(integrityService service.IntegrityService, authenticityService service.AuthenticityService, encryptionService service.EncryptionService,
	availabilityService service.AvailabilityService, fileService service.FileService, notifyService service.NotifyService) *ProcessService {

	return &ProcessService{
		integrityService:    integrityService,
		authenticityService: authenticityService,
		encryptionService:   encryptionService,
		availabilityService: availabilityService,
		fileService:         fileService,
		notifyService:       notifyService,
	}
}

func (s ProcessService) Process(ctx context.Context, req process_request.ProcessRequest) (*process_response.ProcessResponse, error) {
	record, err := s.fileService.GetRecord(ctx, req.Data())
	if err != nil {
		return nil, err
	}

	fileHash, err := record.GetHash()
	if err != nil {
		return nil, err
	}

	certification := domain.Certification{
		Data:   req.Data(),
		Hash:   fileHash,
		Record: record,
	}

	responseBuilder := process_response.NewProcessResponseBuilder()

	if req.IsAuthenticityEnabled() {
		signature, record, err := s.authenticityService.Sign(ctx, authenticity_request.NewSignRequest(
			config.Configuration.AuthenticityPublicKey,
			&config.Configuration.AuthenticityPrivateKey,
			req.AuthenticityKeySource(),
			req.AuthenticityKeyID(),
			req.AuthenticityKeyType(),
			req.AuthenticityUseEnsResolution(),
			req.Data(),
		))
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
		responseBuilder.SignResponse(*authenticity_response.NewSignResponse(signature))
	}

	if req.IsIntegrityEnabled() {
		newCertification, err := s.integrityService.CertifyData(ctx, certification.Data)
		if err != nil {
			return nil, err
		}
		certification = newCertification
		responseBuilder.CertificationResponse(*integrity_response.NewCertificationResponse(certification.Hash, certification.AnchorID))
	}

	if req.IsEncryptionEnabled() {
		req := encryption_request.NewEncryptRequest(
			config.Configuration.EncryptionPublicKey,
			&config.Configuration.EncryptionPrivateKey,
			req.EncryptionKeySource(),
			req.EncryptionKeyID(),
			req.EncryptionKeyType(),
			req.Data(),
		)
		encryptedRecord, err := s.encryptionService.Encrypt(ctx, req)
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
		dataID, err := s.availabilityService.Upload(ctx, certification.Record, req.HostingType())
		if err != nil {
			return nil, err
		}
		certification.DataID = dataID
		if err = s.integrityService.UpdateCertification(ctx, certification); err != nil {
			return nil, err
		}
		responseBuilder.AvailabilityResponse(*availability_response.NewAvailabilityResponse(certification.DataID, req.HostingType()))
	} else {
		if req.IsIntegrityEnabled() {
			if err = s.fileService.SaveFile(ctx, certification.Data, certification.Hash); err != nil {
				return nil, err
			}
		}
	}
	responseBuilder.HashResponse(certification.Hash)

	if certification.AnchorID == 0 {
		if err = s.notifyService.NotifyClient(ctx, []domain.Certification{certification}); err != nil {
			return nil, err
		}
	}

	return responseBuilder.Build(), nil
}
