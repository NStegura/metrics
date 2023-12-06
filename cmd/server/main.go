package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	timeoutShutdown = time.Second * 10
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
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	config := metricsapi.NewConfig()
	err := config.ParseFlags()
	if err != nil {
		return err
	}
	logger, err := configureLogger(config)
	if err != nil {
		return err
	}

	db := repo.New(
		config.StoreInterval,
		config.FileStoragePath,
		config.Restore,
		logger,
	)

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer log.Print("closed DB")
		defer wg.Done()
		<-ctx.Done()

		db.Shutdown()
	}()

	componentsErrs := make(chan error, 1)

	newServer := metricsapi.New(
		config,
		business.New(db, logger),
		logger,
	)
	go func(errs chan<- error) {
		if err = newServer.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("listen and server has failed: %w", err)
		}
	}(componentsErrs)

	go func() {
		err := db.StartBackup()
		if err != nil {
			logger.Warningf("Backup save err, %s", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case err := <-componentsErrs:
		log.Print(err)
		cancelCtx()
	}

	go func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	}()

	return nil
}

func main() {
	if err := runRest(); err != nil {
		log.Fatal(err)
	}
}
