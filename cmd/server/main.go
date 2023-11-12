package main

import (
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"log"
)

func runRest() error {
	config := metricsapi.NewConfig()
	config.ParseFlags()

	newServer := metricsapi.New(
		config,
		business.New(
			repo.New(),
		),
	)
	if err := newServer.Start(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runRest(); err != nil {
		log.Fatal(err)
	}
}
