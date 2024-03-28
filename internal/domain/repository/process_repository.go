package repository

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"github.com/google/uuid"
)

type ProcessRepository interface {
	SaveProcess(ctx context.Context, process domain.Process, isAggregated bool) error
	UpdateAggregatedAnchorID(ctx context.Context, anchorID int) error
	UpdateStatusByAnchorID(ctx context.Context, anchorID int) error

	FindProcessByID(ctx context.Context, id uuid.UUID) (domain.Process, error)
	List(ctx context.Context, pq utils.PaginationQuery) ([]domain.Process, utils.Pagination, error)
}
