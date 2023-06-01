package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
)

type IntegrityRepository interface {
	Certify(ctx context.Context, bytes [][]byte) (certification []domain.Certification, err error)
}
