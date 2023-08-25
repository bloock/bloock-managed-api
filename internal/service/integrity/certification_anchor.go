package integrity

import (
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/integrity/request"
	"context"
)

type CertificationAnchor struct {
	certificationRepository repository.CertificationRepository
	integrityRepository     repository.IntegrityRepository
	notificationRepository  repository.NotificationRepository
}

func NewUpdateAnchorService(certificationRepository repository.CertificationRepository, notificationRepository repository.NotificationRepository, integrityRepository repository.IntegrityRepository) *CertificationAnchor {
	return &CertificationAnchor{certificationRepository: certificationRepository, integrityRepository: integrityRepository, notificationRepository: notificationRepository}
}

func (c CertificationAnchor) UpdateAnchor(ctx context.Context, updateRequest request.UpdateCertificationAnchorRequest) error {
	anchorID := updateRequest.AnchorId
	anchor, err := c.integrityRepository.GetAnchorByID(ctx, anchorID)
	if err != nil {
		return err
	}
	err = c.certificationRepository.UpdateCertificationAnchor(ctx, anchor)
	if err != nil {
		return err
	}

	certifications, err := c.certificationRepository.GetCertificationsByAnchorID(ctx, anchorID)
	if err != nil {
		return err
	}

	for _, crt := range certifications {
		err := c.notificationRepository.NotifyCertification(crt.Hash(), updateRequest.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}
