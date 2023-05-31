package update

import (
	"bloock-managed-api/integrity/domain/repository"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
)

type CertificationAnchor struct {
	certificationRepository repository.CertificationRepository
}

func NewCertificationAnchor(certificationRepository repository.CertificationRepository) *CertificationAnchor {
	return &CertificationAnchor{certificationRepository: certificationRepository}
}

func (c CertificationAnchor) UpdateAnchor(ctx context.Context, anchor integrity.Anchor) error {
	return c.certificationRepository.UpdateCertificationAnchor(ctx, anchor)
}
