package proof

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"github.com/rs/zerolog"
)

var ErrEmptyMessages = errors.New("messages is empty")
var ErrMaxProofMessagesSize = errors.New("too many messages requested for proof")
var ErrMessageNotFound = errors.New("message not found")
var ErrInvalidMessageHash = errors.New("messages do not have a valid format")
var ErrInconsistentMessages = errors.New("some messages do not have the same root")

type GetProof struct {
	messageRepository   repository.MessageAggregatorRepository
	metadataRepository  repository.MetadataRepository
	integrityRepository repository.IntegrityRepository
	maxProofMessageSize int
	apiKey              string
	logger              zerolog.Logger
}

func NewGetProof(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection, maxProofMessageSize int) *GetProof {
	logger := l.With().Caller().Str("component", "get-proof").Logger()

	return &GetProof{
		apiKey:              config.Configuration.Bloock.ApiKey,
		messageRepository:   bloock_repository.NewMessageAggregatorRepository(ctx, l, ent),
		metadataRepository:  bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		integrityRepository: bloock_repository.NewBloockIntegrityRepository(ctx, l),
		maxProofMessageSize: maxProofMessageSize,
		logger:              logger,
	}
}

func (s GetProof) Get(ctx context.Context, hashes []string) (domain.BloockProof, error) {
	if err := s.validateRequest(hashes); err != nil {
		return domain.BloockProof{}, err
	}

	messages, err := s.messageRepository.GetMessagesByHashes(ctx, hashes)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.BloockProof{}, err
	}

	/*messages, err := s.messageRepository.FindMessagesByHashesAndRoot(ctx, hashes, messageDomain.Root)
	if err != nil {
		return domain.BloockProof{}, err
	}*/
	if len(messages) == 0 {
		err = ErrMessageNotFound
		s.logger.Error().Err(err).Msg("")
		return domain.BloockProof{}, err
	}
	if len(messages) != len(hashes) {
		err = ErrInconsistentMessages
		s.logger.Error().Err(err).Msg("")
		return domain.BloockProof{}, err
	}

	bloockProofs := make([]domain.BloockProof, 0)
	for _, mss := range messages {
		bloockProof, err := mss.Proof.ParseToBloockProof(mss.Hash, mss.Root)
		if err != nil {
			s.logger.Error().Err(err).Msg("")
			return domain.BloockProof{}, err
		}
		bloockProofs = append(bloockProofs, bloockProof)
	}

	differentRoots := make([]string, 0)
	temporaryJoinedProofs := make(map[string][]domain.BloockProof, 0)
	for _, bloockProof := range bloockProofs {
		_, ok := temporaryJoinedProofs[bloockProof.Root]
		if !ok {
			temporaryJoinedProofs[bloockProof.Root] = []domain.BloockProof{bloockProof}
			differentRoots = append(differentRoots, bloockProof.Root)
		} else {
			temporaryJoinedProofs[bloockProof.Root] = append(temporaryJoinedProofs[bloockProof.Root], bloockProof)
		}
	}

	joinedProofs := make([]domain.BloockProof, 0)
	for _, tmpProofs := range temporaryJoinedProofs {
		joinedBloockProof, err := domain.JoinBloockMultiProofs(tmpProofs)
		if err != nil {
			s.logger.Error().Err(err).Msg("")
			return domain.BloockProof{}, err
		}
		joinedProofs = append(joinedProofs, joinedBloockProof)
	}

	certificationProof, err := s.integrityRepository.GetProof(ctx, differentRoots, s.apiKey)
	if err != nil {
		return domain.BloockProof{}, err
	}

	/*_, certificationProof, err := s.metadataRepository.GetCertificationByHashAndAnchorID(ctx, messages[0].Root, messages[0].AnchorID)
	if err != nil {
		return domain.BloockProof{}, err
	}*/

	/*joinedBloockProof, err := domain.JoinBloockMultiProofs(bloockProofs)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.BloockProof{}, err
	}*/

	assembledProof, err := certificationProof.AssembleProof(joinedProofs)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.BloockProof{}, err
	}

	return assembledProof, nil
}

func (s GetProof) validateRequest(hash []string) error {
	if len(hash) == 0 {
		err := ErrEmptyMessages
		s.logger.Error().Err(err).Msg("")
		return err
	} else if len(hash) > s.maxProofMessageSize {
		err := ErrMaxProofMessagesSize
		s.logger.Error().Err(err).Msg("")
		return err
	}

	for _, h := range hash {
		if !utils.IsSHA256(h) {
			return ErrInvalidMessageHash
		}
	}
	return nil
}
