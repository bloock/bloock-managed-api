package service

import (
	"bloock-managed-api/internal/service/authenticity/request"
	"bloock-managed-api/internal/service/authenticity/response"
	update_request "bloock-managed-api/internal/service/integrity/request"
	create_response "bloock-managed-api/internal/service/integrity/response"
	"context"
)

type BaseProcessService interface {
	Process(ctx context.Context, req ProcessRequest) (*response.ProcessResponse, error)
}

type AuthenticityService interface {
	Sign(ctx context.Context, SignRequest request.SignRequest) (string, []byte, error)
}

type IntegrityService interface {
	Certify(ctx context.Context, files []byte) ([]create_response.CertificationResponse, error)
}

type AvailabilityService interface {
	UploadHosted(ctx context.Context, data []byte) (string, error)
	UploadIpfs(ctx context.Context, data []byte) (string, error)
}
type CertificateUpdateAnchorService interface {
	UpdateAnchor(ctx context.Context, updateRequest update_request.UpdateCertificationAnchorRequest) error
}
