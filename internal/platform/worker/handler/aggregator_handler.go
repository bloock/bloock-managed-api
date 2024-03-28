package handler

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/platform/cron"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/service/aggregate"
	"github.com/rs/zerolog"
)

func AggregatorHandler(l zerolog.Logger, ent *connection.EntConnection) cron.CronHandler {
	return func(ctx context.Context) error {
		service := aggregate.NewServiceAggregator(ctx, l, ent)

		err := service.Aggregate(ctx)
		return err
	}
}
