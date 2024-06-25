package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/NStegura/metrics/internal/clients/metric"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
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

func startAgent() error {
	config := agent.NewConfig()
	err := config.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	logger, err := configureLogger(config)
	if err != nil {
		return fmt.Errorf("failed to configure logger: %w", err)
	}
	metricsCli, err := metric.New(
		config.HTTPAddr,
		config.MetricCliKey,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to init metric client: %w", err)
	}
	ag := agent.New(config, metricsCli, logger)
	if err = ag.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}
	return nil
}

func printProjectInfo() {
	var s strings.Builder
	s.WriteString("Build version: ")
	if buildVersion != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildVersion))
	} else {
		s.WriteString(fmt.Sprintf("<%s>\n", buildVersion))
	}

	s.WriteString("Build date: ")
	if buildDate != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildDate))
	} else {
		s.WriteString(fmt.Sprintf("<%s>\n", buildDate))
	}

	s.WriteString("Build commit: ")
	if buildCommit != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildCommit))
	} else {
		s.WriteString(fmt.Sprintf("<%s>\n", buildCommit))
	}
	log.Println(s.String())
}

func main() {
	printProjectInfo()
	if err := startAgent(); err != nil {
		log.Fatal(err)
	}
}
