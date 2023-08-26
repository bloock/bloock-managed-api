package service

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/service/authenticity/request"
	"bloock-managed-api/internal/service/authenticity/response"
	"context"
)

type ProcessService struct {
	integrityService    IntegrityService
	authenticityService AuthenticityService
	availabilityService AvailabilityService
}

func NewProcessService(integrityService IntegrityService, authenticityService AuthenticityService,
	availabilityService AvailabilityService) *ProcessService {
	return &ProcessService{
		integrityService:    integrityService,
		authenticityService: authenticityService,
		availabilityService: availabilityService,
	}
}

func (s ProcessService) Process(ctx context.Context, req ProcessRequest) (*response.ProcessResponse, error) {
	responseBuilder := response.NewProcessResponseBuilder()

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

	if req.HostingType() != NONE {
		switch req.HostingType() {
		case IPFS:
			uploadedDataUrl, err := s.availabilityService.UploadIpfs(ctx, req.Data())
			if err != nil {
				return nil, err
			}
			responseBuilder.AvailabilityResponse(uploadedDataUrl)
			break
		case HOSTED:
			uploadedDataUrl, err := s.availabilityService.UploadHosted(ctx, req.Data())
			if err != nil {
				return nil, err
			}
			responseBuilder.AvailabilityResponse(uploadedDataUrl)
			break
		}
	}

	return responseBuilder.Build(), nil
}
