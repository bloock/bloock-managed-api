package repository

import (
	"context"
	"encoding/json"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/message"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type MessageAggregatorRepository struct {
	connection *connection.EntConnection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func NewMessageAggregatorRepository(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) repository.MessageAggregatorRepository {
	logger := l.With().Caller().Str("component", "message-aggregator-repository").Logger()

	return &MessageAggregatorRepository{
		connection: ent,
		dbTimeout:  5 * time.Second,
		logger:     logger,
	}
}

func mapToRawProof(proof domain.MerkleTreeProof) (json.RawMessage, error) {
	proofBytes, err := json.Marshal(proof)
	if err != nil {
		return nil, err
	}
	var rawProof json.RawMessage
	if err = json.Unmarshal(proofBytes, &rawProof); err != nil {
		return nil, err
	}

	return rawProof, nil
}

func mapToMerkleTreeProof(rawProof json.RawMessage) (domain.MerkleTreeProof, error) {
	var proof domain.MerkleTreeProof
	err := json.Unmarshal(rawProof, &proof)
	if err != nil {
		return domain.MerkleTreeProof{}, err
	}

	return proof, nil
}

func mapToMessage(mss *ent.Message) (domain.Message, error) {
	proof, err := mapToMerkleTreeProof(mss.Proof)
	if err != nil {
		return domain.Message{}, err
	}
	return domain.Message{
		Hash:     mss.Message,
		Root:     mss.Root,
		AnchorID: mss.AnchorID,
		Proof:    proof,
	}, nil
}

func (s MessageAggregatorRepository) SaveMessage(ctx context.Context, message domain.Message) error {
	crt := s.connection.DB().
		Message.Create().
		SetMessage(message.Hash)

	if _, err := crt.Save(ctx); err != nil {
		if ent.IsConstraintError(err) && strings.Contains(err.Error(), "duplicate key") {
			return nil
		}
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s MessageAggregatorRepository) GetPendingMessages(ctx context.Context) ([]domain.Message, error) {
	messagesSchema, err := s.connection.DB().Message.Query().
		Where(message.RootEQ("")).All(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return []domain.Message{}, err
	}

	var messages []domain.Message
	for _, ms := range messagesSchema {
		newMss := domain.Message{
			Hash: ms.Message,
		}
		messages = append(messages, newMss)
	}

	return messages, nil
}

func (s MessageAggregatorRepository) GetMessagesByRootAndAnchorID(ctx context.Context, root string, anchorID int) ([]domain.Message, error) {
	messagesSchema, err := s.connection.DB().Message.Query().
		Where(message.RootEQ(root), message.AnchorIDEQ(anchorID)).All(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return []domain.Message{}, err
	}

	var messages []domain.Message
	for _, ms := range messagesSchema {
		newMss := domain.Message{
			Hash: ms.Message,
		}
		messages = append(messages, newMss)
	}

	return messages, nil
}

func (s MessageAggregatorRepository) FindMessagesByHashesAndRoot(ctx context.Context, hash []string, root string) ([]domain.Message, error) {
	messageSchemas, err := s.connection.DB().Message.Query().
		Where(message.MessageIn(hash...), message.RootEQ(root)).Order(ent.Desc(message.FieldAnchorID)).All(ctx)
	if err != nil {
		if ent.IsNotFound(err) && strings.Contains(err.Error(), "not found") {
			return []domain.Message{}, nil
		}
		s.logger.Error().Err(err).Msg("")
		return []domain.Message{}, err
	}

	var messages []domain.Message
	for _, mss := range messageSchemas {
		proof, err := mapToMerkleTreeProof(mss.Proof)
		if err != nil {
			s.logger.Error().Err(err).Msg("")
			return []domain.Message{}, err
		}
		messages = append(messages, domain.Message{
			Hash:     mss.Message,
			AnchorID: mss.AnchorID,
			Root:     mss.Root,
			Proof:    proof,
		})
	}

	return messages, nil
}

func (s MessageAggregatorRepository) GetMessageByHash(ctx context.Context, hash string) (domain.Message, error) {
	messageSchema, err := s.connection.DB().Message.Query().
		Where(message.MessageEQ(hash)).Order(ent.Desc(message.FieldAnchorID)).First(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.Message{}, err
	}

	messageDomain, err := mapToMessage(messageSchema)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return domain.Message{}, err
	}

	return messageDomain, nil
}

func (s MessageAggregatorRepository) ExistRoot(ctx context.Context, root string) (bool, error) {
	existRoot, err := s.connection.DB().Message.Query().
		Where(message.RootEQ(root)).Exist(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return false, err
	}

	return existRoot, nil
}

func (s MessageAggregatorRepository) UpdateMessage(ctx context.Context, m domain.Message) error {
	rawProof, err := mapToRawProof(m.Proof)
	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	if _, err = s.connection.DB().Message.Update().
		SetAnchorID(m.AnchorID).
		SetRoot(m.Root).
		SetProof(rawProof).
		Where(message.MessageEQ(m.Hash), message.RootEQ("")).Save(ctx); err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}
