package retrieve

import (
	"bloock-managed-api/integrity/domain/repository"
	"context"
)

type CertificationByAnchorAndHash struct {
	certificationRepository repository.CertificationRepository
}

func NewCertificationByAnchorAndHash(certificationRepository repository.CertificationRepository) *CertificationByAnchorAndHash {
	return &CertificationByAnchorAndHash{certificationRepository: certificationRepository}
}

func (c CertificationByAnchorAndHash) Retrieve(ctx context.Context, anchorID int, hash string) (*CertificationByAnchorAndHashResponse, error) {
	certification, err := c.certificationRepository.GetCertification(ctx, anchorID, hash)
	if err != nil {
		return &CertificationByAnchorAndHashResponse{}, err
	}

	return NewCertificationByAnchorAndHashResponse(certification.Hashes()[0], certification.Anchor()), nil
}
