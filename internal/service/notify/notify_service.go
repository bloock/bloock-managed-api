package notify

import (
	"bloock-managed-api/internal/domain/repository"
	"context"
	"errors"
)

var (
	ErrSignKeyNotSupported    = errors.New("key type not supported for signing")
	ErrEncryptKeyNotSupported = errors.New("key type not supported for encrypting")
	ErrUnsupportedHosting     = errors.New("unsupported hosting type")
)

type NotifyService struct {
	availabilityRepository repository.AvailabilityRepository
	metadataRepository     repository.MetadataRepository
	notificationRepository repository.NotificationRepository
}

func NewNotifyService(
	availabilityRepository repository.AvailabilityRepository,
	metadataRepository repository.MetadataRepository,
	notificationRepository repository.NotificationRepository,
) *NotifyService {

	return &NotifyService{
		availabilityRepository: availabilityRepository,
		metadataRepository:     metadataRepository,
		notificationRepository: notificationRepository,
	}
}

func (s NotifyService) Notify(ctx context.Context, anchorID int) error {
	certifications, err := s.metadataRepository.GetCertificationsByAnchorID(ctx, anchorID)
	if err != nil {
		return err
	}

	for _, crt := range certifications {
		var fileBytes []byte
		var err error

		if len(crt.Data) != 0 {
			fileBytes = crt.Data
		} else {
			if crt.DataID != "" {
				fileBytes, err = s.availabilityRepository.FindFile(ctx, crt.DataID)
				if err != nil {
					return err
				}
			} else {
				fileBytes, err = s.availabilityRepository.RetrieveTmp(ctx, crt.Hash)
				if err != nil {
					return err
				}
			}
		}

		if err = s.notificationRepository.NotifyCertification(crt.Hash, fileBytes); err != nil {
			return err
		}
	}

	return nil
}
