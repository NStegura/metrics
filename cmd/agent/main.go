package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/NStegura/metrics/config"

	"github.com/NStegura/metrics/internal/clients/metric"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func configureLogger(config *config.AgentConfig) (*logrus.Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	logger.SetLevel(level)
	return logger, nil
}

func startAgent() error {
	cfg := config.NewAgentConfig()
	err := cfg.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	logger, err := configureLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to configure logger: %w", err)
	}
	metricsCli, err := metric.New(
		cfg.HTTPAddr,
		cfg.BodyHashKey,
		cfg.PublicCryptoKey,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to init metric client: %w", err)
	}
	ag := agent.New(cfg, metricsCli, logger)
	if err = ag.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}
	return nil
}

func printProjectInfo() {
	var s strings.Builder
	na := "N/A"
	s.WriteString("Build version: ")
	if buildVersion != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildVersion))
	} else {
		s.WriteString(na)
	}

	s.WriteString("Build date: ")
	if buildDate != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildDate))
	} else {
		s.WriteString(na)
	}

	s.WriteString("Build commit: ")
	if buildCommit != "" {
		s.WriteString(fmt.Sprintf("<%s>\n", buildCommit))
	} else {
		s.WriteString(na)
	}
	log.Println(s.String())
}

func main() {
	printProjectInfo()
	if err := startAgent(); err != nil {
		log.Fatal(err)
	}
}
