package main

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/platform/cron"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/worker"
	"github.com/bloock/bloock-managed-api/pkg"
	"github.com/getsentry/sentry-go"
	"log"
	"sync"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
)

func main() {
	_, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger := pkg.InitLogger(config.Configuration.Api.DebugMode)

	// Initialize tracer if set
	if config.Configuration.Tracing.Enabled {
		options := sentry.ClientOptions{
			Dsn:         config.Configuration.Tracing.TracerConnUrl,
			Environment: config.Configuration.Tracing.AppEnvironment,
			Release:     config.Configuration.Tracing.AppVersion,
		}
		if config.Configuration.Tracing.AppEnvironment == "production" {
			options.EnableTracing = true
			options.TracesSampleRate = 1.0
			options.TracesSampler = func(ctx sentry.SamplingContext) float64 {
				if ctx.Span.Op == "http.server" {
					return 1.0
				}
				return 0.0
			}
		}
		if err = sentry.Init(options); err != nil {
			log.Fatal(err.Error())
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(config.Configuration.Db.ConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}

	go func() {
		defer wg.Done()
		server, err := rest.NewServer(logger, conn, config.Configuration.Integrity.MaxProofMessageSize)
		if err != nil {
			panic(err)
		}
		err = server.Start()
		if err != nil {
			panic(err)
		}
	}()

	if config.Configuration.Integrity.AggregateMode && config.Configuration.Integrity.AggregateWorker {
		ctx := context.Background()
		wg.Add(1)

		c := cron.NewClientCron()

		go func() {
			defer wg.Done()
			cr, err := worker.NewCronWorker(c, logger, config.Configuration.Integrity.AggregateInterval, conn)
			if err != nil {
				panic(err)
			}
			cr.Run(ctx)
		}()
	}

	wg.Wait()
}
