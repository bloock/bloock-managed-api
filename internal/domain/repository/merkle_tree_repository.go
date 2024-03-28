package repository

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
)

type MerkleTreeRepository interface {
	Create(ctx context.Context, messages []domain.Message) (domain.MerkleTree, error)
}
