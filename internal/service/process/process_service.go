package process

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/authenticity/request"
	"bloock-managed-api/internal/service/authenticity/response"
	integrity_response "bloock-managed-api/internal/service/integrity/response"
	process_request "bloock-managed-api/internal/service/process/request"
	process_response "bloock-managed-api/internal/service/process/response"
	"context"
)

type ProcessService struct {
	integrityService    service.IntegrityService
	authenticityService service.AuthenticityService
	availabilityService service.AvailabilityService
	fileService         service.FileService
}

func NewProcessService(integrityService service.IntegrityService, authenticityService service.AuthenticityService,
	availabilityService service.AvailabilityService, fileService service.FileService) *ProcessService {
	return &ProcessService{
		integrityService:    integrityService,
		authenticityService: authenticityService,
		availabilityService: availabilityService,
		fileService:         fileService,
	}
}

func (s ProcessService) Process(ctx context.Context, req process_request.ProcessRequest) (*process_response.ProcessResponse, error) {
	responseBuilder := process_response.NewProcessResponseBuilder()
	var certifications *integrity_response.CertificationResponse
	if req.IsIntegrityEnabled() {
		certifications, err := s.integrityService.Certify(ctx, req.Data())
		if err != nil {
			return nil, err
		}
		responseBuilder.CertificationResponse(certifications)
	}

	if req.IsAuthenticityEnabled() {
		var signature, signedData, err = s.authenticityService.
			Sign(ctx, *request.NewSignRequest(
				config.Configuration.PublicKey,
				&config.Configuration.PrivateKey,
				req.KeyID(),
				req.Kty(),
				req.KeyType(),
				req.Data(),
				req.UseEnsResolution(),
			))
		if err != nil {
			return nil, err
		}

		req.ReplaceDataWith(signedData)
		responseBuilder.SignResponse(*response.NewSignResponse(signature))
	}

	if req.HostingType() != domain.NONE {

		dataID, err := s.availabilityService.Upload(ctx, req.Data(), req.HostingType())
		if err != nil {
			return nil, err
		}
		if err := s.integrityService.SetDataIDToCertification(ctx, certifications.Hash(), dataID); err != nil {
			return nil, err
		}
		responseBuilder.AvailabilityResponse(dataID)
	} else {

		if err := s.fileService.SaveFile(ctx, req.Data()); err != nil {
			return nil, err
		}
	}

	return responseBuilder.Build(), nil
}
