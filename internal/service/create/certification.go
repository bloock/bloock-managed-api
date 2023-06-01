package create

import (
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/response"
	"context"
)

type Certification struct {
	certificationRepository repository.CertificationRepository
	integrityRepository     repository.IntegrityRepository
}

func NewCertification(certificationRepository repository.CertificationRepository, integrityRepository repository.IntegrityRepository) *Certification {
	return &Certification{certificationRepository: certificationRepository, integrityRepository: integrityRepository}
}

func (c Certification) Certify(ctx context.Context, files [][]byte) ([]response.CertificationResponse, error) {
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
