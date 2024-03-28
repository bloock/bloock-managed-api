package worker

import (
	"github.com/bloock/bloock-managed-api/internal/platform/cron"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/worker/handler"
	"github.com/rs/zerolog"
	"time"
)

func NewCronWorker(c *cron.ClientCron, l zerolog.Logger, newAnchorInterval int, ent *connection.EntConnection) (*cron.CronRuntime, error) {
	runtime, err := cron.NewCronRuntime(c, 5*time.Second, l)
	if err != nil {
		return nil, err
	}

	runtime.AddHandler("aggregate-worker", time.Second*time.Duration(newAnchorInterval), handler.AggregatorHandler(l, ent))

	return runtime, nil
}
