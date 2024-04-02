package notify

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/rs/zerolog"
)

type ServiceNotifier struct {
	availabilityRepository repository.AvailabilityRepository
	metadataRepository     repository.MetadataRepository
	notificationRepository repository.NotificationRepository
	integrityRepository    repository.IntegrityRepository
	messageRepository      repository.MessageAggregatorRepository
	processRepository      repository.ProcessRepository
	apiKey                 string
	logger                 zerolog.Logger
}

func NewServiceNotifier(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *ServiceNotifier {
	logger := l.With().Caller().Str("component", "service-notifier").Logger()

	return &ServiceNotifier{
		apiKey:                 config.Configuration.Bloock.ApiKey,
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, l),
		metadataRepository:     bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		notificationRepository: bloock_repository.NewHttpNotificationRepository(ctx, l),
		integrityRepository:    bloock_repository.NewBloockIntegrityRepository(ctx, l),
		messageRepository:      bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		processRepository:      bloock_repository.NewProcessRepository(ctx, l, ent),
		logger:                 logger,
	}
}

func (s ServiceNotifier) Notify(ctx context.Context, anchorID int) ([]domain.Certification, error) {
	aggregateCertifications := make([]domain.Certification, 0)
	certifications, err := s.metadataRepository.GetCertificationsByAnchorID(ctx, anchorID)
	if err != nil {
		return nil, err
	}

	for _, crt := range certifications {
		var fileBytes []byte
		var err error

		if s.apiKey != "" {
			if err = s.saveProof(ctx, crt); err != nil {
				return nil, err
			}
		}

		if len(crt.Data) != 0 {
			fileBytes = crt.Data
		} else {
			if crt.DataID != "" {
				fileBytes, err = s.availabilityRepository.FindFile(ctx, crt.DataID)
				if err != nil {
					fileBytes, err = s.availabilityRepository.RetrieveLocal(ctx, crt.DataID)
					if err != nil {
						return nil, err
					}
				}
			} else {
				fileBytes, err = s.availabilityRepository.RetrieveTmp(ctx, crt.Hash)
				if err != nil {
					exists, err := s.messageRepository.ExistRoot(ctx, crt.Hash)
					if err != nil || !exists {
						continue
					} else {
						aggregateCertifications = append(aggregateCertifications, crt)
						continue
					}
				}
			}
		}

		if err = s.notificationRepository.NotifyCertification(crt.Hash, fileBytes); err != nil {
			return nil, err
		}
	}

	if err = s.processRepository.UpdateStatusByAnchorID(ctx, anchorID); err != nil {
		return nil, err
	}

	return aggregateCertifications, nil
}

func (s ServiceNotifier) saveProof(ctx context.Context, certification domain.Certification) error {
	proof, err := s.integrityRepository.GetProof(ctx, certification.Hash, s.apiKey)
	if err != nil {
		return err
	}

	return s.metadataRepository.UpdateCertificationProof(ctx, certification, proof)
}
