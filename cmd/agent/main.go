package main

import (
	"fmt"
	"log"

	"github.com/NStegura/metrics/internal/clients/metric"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/agent"
)

func configureLogger(config *agent.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	logger.SetLevel(level)
	return logger, nil
}

func main() {
	config := agent.NewConfig()
	err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := configureLogger(config)
	if err != nil {
		log.Fatal(err)
	}
	metricsCli, err := metric.New(
		config.HTTPAddr,
		config.MetricCliKey,
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}
	ag := agent.New(config, metricsCli, logger)
	if err = ag.Start(); err != nil {
		logger.Fatal(err)
	}
}
