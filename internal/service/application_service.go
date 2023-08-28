package service

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service/authenticity/request"
	update_request "bloock-managed-api/internal/service/integrity/request"
	create_response "bloock-managed-api/internal/service/integrity/response"
	request2 "bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"context"
)

type BaseProcessService interface {
	Process(ctx context.Context, req request2.ProcessRequest) (*response.ProcessResponse, error)
}

type AuthenticityService interface {
	Sign(ctx context.Context, SignRequest request.SignRequest) (string, []byte, error)
}

type IntegrityService interface {
	Certify(ctx context.Context, files []byte) ([]create_response.CertificationResponse, error)
}

type AvailabilityService interface {
	Upload(ctx context.Context, data []byte, hostingType domain.HostingType) (string, error)
}

type CertificateUpdateAnchorService interface {
	UpdateAnchor(ctx context.Context, updateRequest update_request.UpdateCertificationAnchorRequest) error
}
