package integrity

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
)

type IntegrityService struct {
	certificationRepository repository.CertificationRepository
	integrityRepository     repository.IntegrityRepository
}

func NewIntegrityService(cr repository.CertificationRepository, ir repository.IntegrityRepository) *IntegrityService {
	return &IntegrityService{
		certificationRepository: cr,
		integrityRepository:     ir,
	}
}

func (c IntegrityService) CertifyData(ctx context.Context, data []byte) (domain.Certification, error) {
	certification, err := c.integrityRepository.Certify(ctx, data)
	if err != nil {
		return domain.Certification{}, err
	}

	if err = c.certificationRepository.SaveCertification(ctx, certification); err != nil {
		return domain.Certification{}, err
	}

	return certification, nil
}

func (c IntegrityService) UpdateCertification(ctx context.Context, certification domain.Certification) error {
	exist, err := c.certificationRepository.ExistCertificationByHash(ctx, certification.Hash)
	if err != nil {
		return err
	}
	if !exist {
		return c.certificationRepository.SaveCertification(ctx, certification)
	} else {
		return c.certificationRepository.UpdateCertificationDataID(ctx, certification)
	}
}
