package main

import (
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"log"
)

func runRest() {
	config := metricsapi.NewConfig()
	config.ParseFlags()

	newServer := metricsapi.New(
		config,
		business.New(
			repo.New(),
		),
	)
	if err := newServer.Start(); err != nil {
		log.Fatal(err)
	}

}
func main() {
	runRest()
}
