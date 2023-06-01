package repository

import (
	"bloock-managed-api/integrity/domain"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
)

type CertificationRepository interface {
	SaveCertification(ctx context.Context, certification []domain.Certification) error
	GetCertificationsByAnchorID(ctx context.Context, anchor int) (certification []domain.Certification, err error)
	UpdateCertificationAnchor(ctx context.Context, anchor integrity.Anchor) error
}
