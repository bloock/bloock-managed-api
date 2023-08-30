package integrity

import (
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/integrity/response"
	"context"
)

type IntegrityService struct {
	certificationRepository repository.CertificationRepository
	integrityRepository     repository.IntegrityRepository
}

func NewIntegrityService(certificationRepository repository.CertificationRepository, integrityRepository repository.IntegrityRepository) *IntegrityService {
	return &IntegrityService{certificationRepository: certificationRepository, integrityRepository: integrityRepository}
}
func (c IntegrityService) Certify(ctx context.Context, file []byte) (response.CertificationResponse, error) {
	certification, err := c.integrityRepository.Certify(ctx, file)
	if err != nil {
		return response.CertificationResponse{}, err
	}

	if err := c.certificationRepository.SaveCertification(ctx, certification); err != nil {
		return response.CertificationResponse{}, err
	}

	return *response.NewCertificationResponse(certification.Hash(), certification.AnchorID()), nil
}

func (c IntegrityService) SetDataIDToCertification(ctx context.Context, hash string, id string) error {
	return c.certificationRepository.UpdateCertificationDataID(ctx, hash, id)
}
