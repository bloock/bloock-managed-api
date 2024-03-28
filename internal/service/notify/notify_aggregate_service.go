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

type ServiceAggregateNotifier struct {
	availabilityRepository repository.AvailabilityRepository
	notificationRepository repository.NotificationRepository
	messageRepository      repository.MessageAggregatorRepository
	clientEndpointURL      string
	logger                 zerolog.Logger
}

func NewServiceAggregateNotifier(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *ServiceAggregateNotifier {
	logger := l.With().Caller().Str("component", "service-aggregator-notifier").Logger()

	return &ServiceAggregateNotifier{
		clientEndpointURL:      config.Configuration.Webhook.ClientEndpointUrl,
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, l),
		notificationRepository: bloock_repository.NewHttpNotificationRepository(ctx, l),
		messageRepository:      bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		logger:                 logger,
	}
}

func (s ServiceAggregateNotifier) Notify(ctx context.Context, certifications []domain.Certification) error {
	if s.clientEndpointURL == "" {
		return nil
	}
	for _, crt := range certifications {
		messages, err := s.messageRepository.GetMessagesByRootAndAnchorID(ctx, crt.Hash, crt.AnchorID)
		if err != nil {
			return err
		}

		for _, m := range messages {
			fileBytes, err := s.availabilityRepository.RetrieveTmp(ctx, m.Hash)
			if err != nil {
				return err
			}

			if err = s.notificationRepository.NotifyCertification(m.Hash, fileBytes); err != nil {
				continue
			}
		}
	}
	return nil
}
