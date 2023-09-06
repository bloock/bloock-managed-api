package notify

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/domain/repository"
	"context"
)

type NotifyService struct {
	notificationRepository repository.NotificationRepository
	localStorageRepository repository.LocalStorageRepository
	availabilityRepository repository.AvailabilityRepository
}

func NewNotifyService(nr repository.NotificationRepository, lsr repository.LocalStorageRepository, ar repository.AvailabilityRepository) *NotifyService {
	return &NotifyService{
		notificationRepository: nr,
		localStorageRepository: lsr,
		availabilityRepository: ar,
	}
}

func (n NotifyService) NotifyClient(ctx context.Context, certifications []domain.Certification) error {
	for _, crt := range certifications {
		var fileBytes []byte
		var err error

		if len(crt.Data) != 0 {
			fileBytes = crt.Data
		} else {
			if crt.DataID != "" {
				fileBytes, err = n.availabilityRepository.FindFile(ctx, crt.DataID)
				if err != nil {
					return err
				}
			} else {
				fileBytes, err = n.localStorageRepository.Retrieve(ctx, config.Configuration.FileDir, crt.Hash)
				if err != nil {
					return err
				}
			}
		}

		if err = n.notificationRepository.NotifyCertification(crt.Hash, fileBytes); err != nil {
			return err
		}
	}

	return nil
}
