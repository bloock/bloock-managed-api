package update

import (
	"bloock-managed-api/integrity/domain/repository"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
)

type CertificationAnchor struct {
	certificationRepository repository.CertificationRepository
	notificationRepository  repository.NotificationRepository
}

func NewCertificationAnchor(certificationRepository repository.CertificationRepository, notificationRepository repository.NotificationRepository) *CertificationAnchor {
	return &CertificationAnchor{certificationRepository: certificationRepository, notificationRepository: notificationRepository}
}

func (c CertificationAnchor) UpdateAnchor(ctx context.Context, anchor integrity.Anchor) error {
	err := c.certificationRepository.UpdateCertificationAnchor(ctx, anchor)
	if err != nil {
		return err
	}

	certifications, err := c.certificationRepository.GetCertificationsByAnchorID(ctx, int(anchor.Id))
	if err != nil {
		return err
	}

	for _, crt := range certifications {
		err := c.notificationRepository.NotifyCertification(crt.Hash(), *crt.Anchor())
		if err != nil {
			return err
		}
	}
	return nil
}
