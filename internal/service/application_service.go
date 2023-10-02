package service

import (
	"bloock-managed-api/internal/domain"
	authenticityRequest "bloock-managed-api/internal/service/authenticity/request"
	encryptionRequest "bloock-managed-api/internal/service/encryption/request"
	processRequest "bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type BaseProcessService interface {
	Process(ctx context.Context, req processRequest.ProcessRequest) (*response.ProcessResponse, error)
}

type AuthenticityService interface {
	Sign(ctx context.Context, signRequest authenticityRequest.SignRequest) (string, *record.Record, error)
}

type EncryptionService interface {
	Encrypt(ctx context.Context, request encryptionRequest.EncryptRequest) (*record.Record, error)
}

type IntegrityService interface {
	CertifyData(ctx context.Context, data []byte) (domain.Certification, error)
	UpdateCertification(ctx context.Context, certification domain.Certification) error
}

type AvailabilityService interface {
	Upload(ctx context.Context, record *record.Record, hostingType domain.HostingType) (string, error)
	Download(ctx context.Context, url string) ([]byte, error)
}

type CertificateUpdateAnchorService interface {
	GetCertificationsByAnchorID(ctx context.Context, anchorID int) ([]domain.Certification, error)
}

type FileService interface {
	GetRecord(ctx context.Context, file []byte) (*record.Record, error)
	GetFileHash(ctx context.Context, file []byte) (string, error)
	SaveFile(ctx context.Context, file []byte, hash string) error
}

type NotifyService interface {
	NotifyClient(ctx context.Context, certifications []domain.Certification) error
}
