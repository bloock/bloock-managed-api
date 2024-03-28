package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"time"
)

type ClientCron struct {
	scheduler *gocron.Scheduler
	handlers  []cronJob
}

type CronHandler func(context.Context) error
type cronJob struct {
	ctx     context.Context
	name    string
	spec    time.Duration
	fixTime string
	job     CronHandler
	logger  zerolog.Logger
}

func newCronJob(name string, spec time.Duration, fixTime string, job CronHandler) cronJob {
	return cronJob{
		name:    name,
		spec:    spec,
		fixTime: fixTime,
		job:     job,
	}
}

func NewClientCron() *ClientCron {
	c := gocron.NewScheduler(time.UTC)

	return &ClientCron{
		scheduler: c,
		handlers:  make([]cronJob, 0),
	}
}

func (a *ClientCron) AddJob(name string, spec time.Duration, fixTime string, handler CronHandler) {
	job := newCronJob(name, spec, fixTime, handler)
	a.handlers = append(a.handlers, job)
}

func (a *ClientCron) Start(ctx context.Context) error {
	for _, handler := range a.handlers {
		if handler.fixTime != "" {
			_, err := a.scheduler.Cron(handler.fixTime).Do(handler.Run)
			if err != nil {
				return err
			}
			continue
		}
		_, err := a.scheduler.Every(handler.spec).Do(handler.Run)
		if err != nil {
			return err
		}
	}

	a.scheduler.StartAsync()

	return nil
}

func (a *ClientCron) Close(shutdownTime time.Duration) error {
	stop := make(chan bool)
	go func() {
		a.scheduler.Stop()
		stop <- true
	}()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	select {
	case <-stop:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("couldn't close cron client before timeout")
	}
}

func (c cronJob) Run() {
	ctx := context.Background()

	err := c.job(ctx)
	if err != nil {
		c.logger.Error().Str("job-name", c.name).Msgf("error running cron: %s", err.Error())
		return
	}
	c.logger.Info().Str("job-name", c.name).Msg("job runned successfully")
}
