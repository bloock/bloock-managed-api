package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
)

type IntegrityRepository interface {
	Certify(ctx context.Context, file []byte) (domain.Certification, error)
}
