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
func (c IntegrityService) Certify(ctx context.Context, files []byte) ([]response.CertificationResponse, error) {
	certifications, err := c.integrityRepository.Certify(ctx, files)
	if err != nil {
		return []response.CertificationResponse{}, err
	}

	if err := c.certificationRepository.SaveCertification(ctx, certifications); err != nil {
		return []response.CertificationResponse{}, err
	}

	var responses []response.CertificationResponse
	for _, crt := range certifications {
		responses = append(responses, *response.NewCertificationResponse(crt.Hash(), crt.AnchorID()))
	}

	return responses, nil
}
