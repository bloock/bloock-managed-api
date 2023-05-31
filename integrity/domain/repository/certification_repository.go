package repository

import (
	"bloock-managed-api/integrity/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
)

type CertificationRepository interface {
	SaveCertification(ctx context.Context, certification domain.Certification) error
	GetCertification(ctx context.Context, anchor int, hash string) (certification *domain.Certification, err error)
	UpdateCertificationAnchor(ctx context.Context, anchor integrity.Anchor) error
}
