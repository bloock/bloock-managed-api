package service

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service/authenticity/request"
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
	CertifyData(ctx context.Context, data []byte) (domain.Certification, error)
	UpdateCertification(ctx context.Context, certification domain.Certification) error
}

type AvailabilityService interface {
	Upload(ctx context.Context, data []byte, hostingType domain.HostingType) (string, error)
}

type CertificateUpdateAnchorService interface {
	GetCertificationsByAnchorID(ctx context.Context, anchorID int) ([]domain.Certification, error)
}

type FileService interface {
	GetFileHash(ctx context.Context, file []byte) (string, error)
	SaveFile(ctx context.Context, file []byte, hash string) error
}

type NotifyService interface {
	NotifyClient(ctx context.Context, certifications []domain.Certification) error
}
