package main

import (
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/bll"
	"github.com/NStegura/metrics/internal/dal"
	"log"
)

func runRest() {
	config := metricsapi.NewConfig()
	config.ParseFlags()

	newServer := metricsapi.New(
		config,
		bll.New(
			dal.New(),
		),
	)
	if err := newServer.Start(); err != nil {
		log.Fatal(err)
	}

}
func main() {
	runRest()
}
