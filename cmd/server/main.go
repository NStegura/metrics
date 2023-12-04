package main

import (
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func configureLogger(config *metricsapi.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)
	return logger, nil
}

func runRest() error {
	config := metricsapi.NewConfig()
	err := config.ParseFlags()
	if err != nil {
		return err
	}
	logger, err := configureLogger(config)
	if err != nil {
		return err
	}
	r := repo.New(
		config.StoreInterval,
		config.FileStoragePath,
		config.Restore,
		logger,
	)
	defer r.Shutdown()

	newServer := metricsapi.New(
		config,
		business.New(r, logger),
		logger,
	)

	go func() {
		err := r.StartBackup()
		if err != nil {
			logger.Warningf("Backup save err, %s", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		r.Shutdown()
		logger.Fatal(sig)
	}()

	if err = newServer.Start(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := runRest(); err != nil {
		log.Fatal(err)
	}
}
