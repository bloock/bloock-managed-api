package repository

import (
	"bloock-managed-api/internal/domain"
	"context"
)

type CertificationRepository interface {
	SaveCertification(ctx context.Context, certification domain.Certification) error

	GetCertificationsByAnchorID(ctx context.Context, anchor int) (certification []domain.Certification, err error)
	ExistCertificationByHash(ctx context.Context, hash string) (bool, error)

	UpdateCertificationDataID(ctx context.Context, certification domain.Certification) error
}
