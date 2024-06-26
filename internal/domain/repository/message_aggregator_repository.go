package repository

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
)

type MessageAggregatorRepository interface {
	SaveMessage(ctx context.Context, message domain.Message) error
	UpdateMessage(ctx context.Context, m domain.Message) error

	GetPendingMessages(ctx context.Context) ([]domain.Message, error)
	GetMessagesByRootAndAnchorID(ctx context.Context, root string, anchorID int) ([]domain.Message, error)
	GetMessagesByHashes(ctx context.Context, hash []string) ([]domain.Message, error)
	ExistRoot(ctx context.Context, root string) (bool, error)
}
