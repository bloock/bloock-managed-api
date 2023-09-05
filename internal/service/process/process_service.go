package process

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/authenticity/request"
	authenticity_response "bloock-managed-api/internal/service/authenticity/response"
	availability_response "bloock-managed-api/internal/service/availability/response"
	"bloock-managed-api/internal/service/integrity/response"
	process_request "bloock-managed-api/internal/service/process/request"
	process_response "bloock-managed-api/internal/service/process/response"
	"context"
)

type ProcessService struct {
	integrityService    service.IntegrityService
	authenticityService service.AuthenticityService
	availabilityService service.AvailabilityService
	fileService         service.FileService
	notifyService       service.NotifyService
}

func NewProcessService(integrityService service.IntegrityService, authenticityService service.AuthenticityService,
	availabilityService service.AvailabilityService, fileService service.FileService, notifyService service.NotifyService) *ProcessService {

	return &ProcessService{
		integrityService:    integrityService,
		authenticityService: authenticityService,
		availabilityService: availabilityService,
		fileService:         fileService,
		notifyService:       notifyService,
	}
}

func (s ProcessService) Process(ctx context.Context, req process_request.ProcessRequest) (*process_response.ProcessResponse, error) {
	responseBuilder := process_response.NewProcessResponseBuilder()
	asyncClientResponse := false

	fileHash, err := s.fileService.GetFileHash(ctx, req.Data())
	if err != nil {
		return nil, err
	}

	certification := domain.Certification{
		Data: req.Data(),
		Hash: fileHash,
	}

	if req.IsAuthenticityEnabled() {
		signature, signedData, err := s.authenticityService.Sign(ctx, *request.NewSignRequest(
			config.Configuration.PublicKey,
			&config.Configuration.PrivateKey,
			req.KeySource(),
			req.KeyID(),
			req.KeyType(),
			req.Data(),
			req.UseEnsResolution(),
		))
		if err != nil {
			return nil, err
		}

		certification.Data = signedData
		responseBuilder.SignResponse(*authenticity_response.NewSignResponse(signature))
	}

	if req.IsIntegrityEnabled() {
		asyncClientResponse = true
		newCertification, err := s.integrityService.CertifyData(ctx, certification.Data)
		if err != nil {
			return nil, err
		}
		certification.AnchorID = newCertification.AnchorID
		certification.Hash = newCertification.Hash
		responseBuilder.CertificationResponse(*response.NewCertificationResponse(certification.Hash, certification.AnchorID))
	}

	if req.HostingType() != domain.NONE {
		dataID, err := s.availabilityService.Upload(ctx, certification.Data, req.HostingType())
		if err != nil {
			return nil, err
		}
		certification.DataID = dataID
		if err = s.integrityService.UpdateCertification(ctx, certification); err != nil {
			return nil, err
		}
		responseBuilder.AvailabilityResponse(*availability_response.NewAvailabilityResponse(certification.DataID))
	} else {
		if err = s.fileService.SaveFile(ctx, certification.Data, certification.Hash); err != nil {
			return nil, err
		}
	}

	if !asyncClientResponse {
		if err = s.notifyService.NotifyClient(ctx, []domain.Certification{certification}); err != nil {
			return nil, err
		}
	}

	return responseBuilder.Build(), nil
}
