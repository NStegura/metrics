package main

import (
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"github.com/sirupsen/logrus"
	"log"
)

func configureLogger(config *metricsapi.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)
	return logger, nil
}

func runRest() error {
	config := metricsapi.NewConfig()
	config.ParseFlags()
	logger, err := configureLogger(config)
	if err != nil {
		return err
	}

	newServer := metricsapi.New(
		config,
		business.New(
			repo.New(logger),
			logger,
		),
		logger,
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
