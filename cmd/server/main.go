package main

import (
	"context"
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/bll"
	"github.com/NStegura/metrics/internal/dal"
	"log"
)

func runRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	newServer := metricsapi.New(
		metricsapi.NewConfig(),
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
