package process

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/authenticity/request"
	"bloock-managed-api/internal/service/authenticity/response"
	request2 "bloock-managed-api/internal/service/process/request"
	response2 "bloock-managed-api/internal/service/process/response"
	"context"
)

type ProcessService struct {
	integrityService    service.IntegrityService
	authenticityService service.AuthenticityService
	availabilityService service.AvailabilityService
}

func NewProcessService(integrityService service.IntegrityService, authenticityService service.AuthenticityService,
	availabilityService service.AvailabilityService) *ProcessService {
	return &ProcessService{
		integrityService:    integrityService,
		authenticityService: authenticityService,
		availabilityService: availabilityService,
	}
}

func (s ProcessService) Process(ctx context.Context, req request2.ProcessRequest) (*response2.ProcessResponse, error) {
	responseBuilder := response2.NewProcessResponseBuilder()

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
				config.Configuration.PrivateKey,
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
		responseBuilder.AvailabilityResponse(dataID)
	}

	return responseBuilder.Build(), nil
}
