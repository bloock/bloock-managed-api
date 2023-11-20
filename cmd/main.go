package main

import (
	"os"
	"sync"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	_, err := config.InitConfig(logger)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		server, err := rest.NewServer(logger)
		if err != nil {
			panic(err)
		}
		err = server.Start()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
