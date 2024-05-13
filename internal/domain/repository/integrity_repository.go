package repository

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
)

type IntegrityRepository interface {
	Certify(ctx context.Context, file []byte) (domain.Certification, error)
	CertifyFromHash(ctx context.Context, hash string, apiKey string) (domain.Certification, error)
	GetProof(ctx context.Context, hash []string, apiKey string) (domain.BloockProof, error)
}
