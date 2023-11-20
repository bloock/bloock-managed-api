package notify

import (
	"context"
	"errors"

	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/rs/zerolog"
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
	logger                 zerolog.Logger
}

func NewNotifyService(ctx context.Context, l zerolog.Logger) *NotifyService {

	return &NotifyService{
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, l),
		metadataRepository:     bloock_repository.NewBloockMetadataRepository(ctx, l),
		notificationRepository: bloock_repository.NewHttpNotificationRepository(ctx, l),
		logger:                 l,
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
