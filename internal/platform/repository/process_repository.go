package repository

import (
	"context"
	"encoding/json"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/process"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	pkg "github.com/bloock/bloock-managed-api/pkg/response"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type ProcessRepository struct {
	connection *connection.EntConnection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func NewProcessRepository(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) repository.ProcessRepository {
	logger := l.With().Caller().Str("component", "process-repository").Logger()

	return &ProcessRepository{
		connection: ent,
		dbTimeout:  5 * time.Second,
		logger:     logger,
	}
}

func mapToRawProcessResponse(process domain.Process) (json.RawMessage, error) {
	processBytes, err := json.Marshal(process.ProcessResponse)
	if err != nil {
		return nil, err
	}

	return processBytes, nil
}

func mapToProcess(pr *ent.Process) (domain.Process, error) {
	var processResponse pkg.ProcessResponse
	if err := json.Unmarshal(pr.ProcessResponse, &processResponse); err != nil {
		return domain.Process{}, err
	}

	if processResponse.Integrity != nil {
		processResponse.Integrity.AnchorId = pr.AnchorID
	}

	return domain.Process{
		ID:              pr.ID.String(),
		Filename:        pr.Filename,
		Status:          pr.Status,
		CreatedAt:       pr.CreatedAt,
		ProcessResponse: processResponse,
	}, nil
}

func (p ProcessRepository) SaveProcess(ctx context.Context, process domain.Process, isAggregated bool) error {
	processResponse, err := mapToRawProcessResponse(process)
	if err != nil {
		p.logger.Error().Err(err).Msg("")
		return err
	}

	proc := p.connection.DB().
		Process.Create().
		SetFilename(process.Filename).
		SetHash(process.ProcessResponse.Hash).
		SetProcessResponse(processResponse).
		SetAnchorID(process.ProcessResponse.Integrity.AnchorId).
		SetIsAggregated(isAggregated)

	if _, err = proc.Save(ctx); err != nil {
		p.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (p ProcessRepository) FindProcessByID(ctx context.Context, id uuid.UUID) (domain.Process, error) {
	processSchema, err := p.connection.DB().Process.Query().
		Where(process.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) && strings.Contains(err.Error(), "not found") {
			return domain.Process{}, nil
		}
		p.logger.Error().Err(err).Msg("")
		return domain.Process{}, err
	}

	pr, err := mapToProcess(processSchema)
	if err != nil {
		p.logger.Error().Err(err).Msg("")
		return domain.Process{}, err
	}

	return pr, nil
}

func (p ProcessRepository) List(ctx context.Context, pq utils.PaginationQuery) ([]domain.Process, utils.Pagination, error) {
	processSchemas, err := p.connection.DB().Process.Query().
		Order(ent.Desc(process.FieldCreatedAt)).Limit(pq.PerPage).Offset(pq.Skip()).All(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("")
		return []domain.Process{}, utils.Pagination{}, err
	}

	var processes []domain.Process
	for _, pr := range processSchemas {
		prs, err := mapToProcess(pr)
		if err != nil {
			return []domain.Process{}, utils.Pagination{}, err
		}
		processes = append(processes, prs)
	}

	total, err := p.connection.DB().Process.Query().Count(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("")
		return []domain.Process{}, utils.Pagination{}, err
	}

	return processes, utils.NewPagination(pq.Page, pq.PerPage, total), nil
}

func (p ProcessRepository) UpdateAggregatedAnchorID(ctx context.Context, anchorID int) error {
	if _, err := p.connection.DB().Process.Update().
		SetAnchorID(anchorID).
		Where(process.IsAggregatedEQ(true), process.AnchorIDEQ(0)).Save(ctx); err != nil {
		p.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (p ProcessRepository) UpdateStatusByAnchorID(ctx context.Context, anchorID int) error {
	if _, err := p.connection.DB().Process.Update().
		SetStatus(true).
		Where(process.AnchorIDEQ(anchorID)).Save(ctx); err != nil {
		p.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}
