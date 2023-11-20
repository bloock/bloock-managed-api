package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain"
)

type IntegrityRepository interface {
	Certify(ctx context.Context, file []byte) (domain.Certification, error)
}
