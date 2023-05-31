package create

import (
	"bloock-managed-api/integrity/domain/repository"
	"context"
)

type Certification struct {
	certificationRepository repository.CertificationRepository
	integrityRepository     repository.IntegrityRepository
}

func NewCertification(certificationRepository repository.CertificationRepository, integrityRepository repository.IntegrityRepository) *Certification {
	return &Certification{certificationRepository: certificationRepository, integrityRepository: integrityRepository}
}

func (c Certification) Certify(ctx context.Context, files [][]byte) error {
	certification, err := c.integrityRepository.Certify(ctx, files)
	if err != nil {
		return err
	}

	if err := c.certificationRepository.SaveCertification(ctx, certification); err != nil {
		return err
	}

	return nil
}
