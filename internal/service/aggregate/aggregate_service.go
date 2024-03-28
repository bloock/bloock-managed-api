package aggregate

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/rs/zerolog"
)

var (
	ErrApiKeyNotFound         = errors.New("api key not found")
	ErrMinimumPendingMessages = errors.New("minimum number of messages to aggregate is 2")
)

type ServiceAggregator struct {
	integrityRepository  repository.IntegrityRepository
	merkleTreeRepository repository.MerkleTreeRepository
	messageRepository    repository.MessageAggregatorRepository
	metadataRepository   repository.MetadataRepository
	processRepository    repository.ProcessRepository
	apiKey               string
	logger               zerolog.Logger
}

func NewServiceAggregator(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *ServiceAggregator {
	logger := l.With().Caller().Str("component", "service-aggregator").Logger()

	apiKey := pkg.GetApiKeyFromContext(ctx)
	if apiKey == "" {
		apiKey = config.Configuration.Bloock.ApiKey
	}
	return &ServiceAggregator{
		apiKey:               apiKey,
		integrityRepository:  bloock_repository.NewBloockIntegrityRepository(ctx, l),
		merkleTreeRepository: bloock_repository.NewMerkleTreeRepository(l),
		messageRepository:    bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		metadataRepository:   bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		processRepository:    bloock_repository.NewProcessRepository(ctx, l, ent),
		logger:               logger,
	}
}

func (s ServiceAggregator) Aggregate(ctx context.Context) error {
	if s.apiKey == "" {
		return ErrApiKeyNotFound
	}
	pendingMessages, err := s.messageRepository.GetPendingMessages(ctx)
	if err != nil {
		return err
	}
	if len(pendingMessages) <= 1 {
		err = ErrMinimumPendingMessages
		s.logger.Error().Err(err).Msg("")
		return err
	}

	merkleTree, err := s.merkleTreeRepository.Create(ctx, pendingMessages)
	if err != nil {
		return err
	}

	rootCertification, err := s.integrityRepository.CertifyFromHash(ctx, merkleTree.Root, s.apiKey)
	if err != nil {
		return err
	}
	if err = s.metadataRepository.SaveCertification(ctx, rootCertification); err != nil {
		return err
	}

	for _, pm := range pendingMessages {
		messageProof, ok := merkleTree.Proof[pm.Hash]
		if !ok {
			continue
		}
		updatedMessage := domain.Message{
			Hash:     pm.Hash,
			Proof:    messageProof,
			AnchorID: rootCertification.AnchorID,
			Root:     rootCertification.Hash,
		}
		if err = s.messageRepository.UpdateMessage(ctx, updatedMessage); err != nil {
			s.logger.Error().Err(err).Msg("")
			continue
		}
	}

	return s.processRepository.UpdateAggregatedAnchorID(ctx, rootCertification.AnchorID)
}
