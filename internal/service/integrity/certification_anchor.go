package integrity

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
)

type CertificationAnchor struct {
	certificationRepository repository.CertificationRepository
}

func NewUpdateAnchorService(certificationRepository repository.CertificationRepository) *CertificationAnchor {
	return &CertificationAnchor{
		certificationRepository: certificationRepository,
	}
}

func (c CertificationAnchor) GetCertificationsByAnchorID(ctx context.Context, anchorID int) ([]domain.Certification, error) {
	certifications, err := c.certificationRepository.GetCertificationsByAnchorID(ctx, anchorID)
	if err != nil {
		return []domain.Certification{}, err
	}

	return certifications, nil
}
