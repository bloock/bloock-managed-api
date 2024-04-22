package repository

import (
	"context"
	"encoding/json"

	"github.com/bloock/bloock-managed-api/internal/domain"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type MetadataRepository interface {
	SaveCertification(ctx context.Context, certification domain.Certification) error
	GetCertificationsByAnchorID(ctx context.Context, anchor int) (certification []domain.Certification, err error)
	GetCertificationByHashAndAnchorID(ctx context.Context, hash string, anchorID int) (domain.Certification, domain.BloockProof, error)
	FindCertificationByHash(ctx context.Context, hash string) (domain.Certification, error)
	ExistCertificationByHash(ctx context.Context, hash string) (bool, error)
	UpdateCertificationDataID(ctx context.Context, certification domain.Certification) error
	UpdateCertification(ctx context.Context, certification domain.Certification) error
	UpdateCertificationProof(ctx context.Context, cert domain.Certification, proof json.RawMessage) error

	GetRecord(ctx context.Context, file []byte) (*record.Record, error)
	GetRecordDetails(ctx context.Context, file []byte) (*record.RecordDetails, error)
	GetFileHash(ctx context.Context, file []byte) (string, error)
}
