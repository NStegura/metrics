package main

import (
	"github.com/NStegura/metrics/internal/app/agent"
	"github.com/sirupsen/logrus"
	"log"
)

func configureLogger(config *agent.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
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

	ag := agent.New(config, logger)
	if err = ag.Start(); err != nil {
		log.Fatal(err)
	}
}
