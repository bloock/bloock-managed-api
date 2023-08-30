package integrity

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain/repository"
	"bloock-managed-api/internal/service/integrity/request"
	"context"
)

type CertificationAnchor struct {
	certificationRepository repository.CertificationRepository
	localStorageRepository  repository.LocalStorageRepository
	availabilityRepository  repository.AvailabilityRepository
	integrityRepository     repository.IntegrityRepository
	notificationRepository  repository.NotificationRepository
}

func NewUpdateAnchorService(
	certificationRepository repository.CertificationRepository,
	notificationRepository repository.NotificationRepository,
	integrityRepository repository.IntegrityRepository,
	availabilityRepository repository.AvailabilityRepository,
	storageRepository repository.LocalStorageRepository,
) *CertificationAnchor {
	return &CertificationAnchor{
		certificationRepository: certificationRepository,
		localStorageRepository:  storageRepository,
		availabilityRepository:  availabilityRepository,
		integrityRepository:     integrityRepository,
		notificationRepository:  notificationRepository,
	}
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
		var fileBytes []byte
		fileBytes, err = c.availabilityRepository.FindFile(ctx, crt.DataID())
		if err != nil {
			return err
		}
		if len(fileBytes) == 0 {
			fileBytes, err = c.localStorageRepository.Retrieve(ctx, config.Configuration.FileDir, crt.Hash())
			if err != nil {
				return err
			}
		}
		err := c.notificationRepository.NotifyCertification(crt.Hash(), updateRequest.Payload, fileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
