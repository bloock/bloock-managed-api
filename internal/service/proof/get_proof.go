package proof

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/rs/zerolog"
)

var ErrMessageNotFound = errors.New("message not found")

type GetProof struct {
	messageRepository  repository.MessageAggregatorRepository
	metadataRepository repository.MetadataRepository
	logger             zerolog.Logger
}

func NewGetProof(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *GetProof {
	logger := l.With().Caller().Str("component", "get-proof").Logger()

	return &GetProof{
		messageRepository:  bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		metadataRepository: bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		logger:             logger,
	}
}

func (s GetProof) Get(ctx context.Context, hash string) (domain.Proof, error) {
	message, err := s.messageRepository.FindMessageByHash(ctx, hash)
	if err != nil {
		return domain.Proof{}, err
	}
	if message.Root == "" {
		err = ErrMessageNotFound
		s.logger.Error().Err(err).Msg("")
		return domain.Proof{}, err
	}

	_, certificationProof, err := s.metadataRepository.GetCertificationByHashAndAnchorID(ctx, message.Root, message.AnchorID)
	if err != nil {
		return domain.Proof{}, err
	}

	convertedBloockProof, err := message.Proof.ConvertToBloockProof(message.Hash)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.Proof{}, err
	}

	assembledProof, err := certificationProof.AssembleProof([]domain.Proof{convertedBloockProof})
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.Proof{}, err
	}

	return assembledProof, nil
}
