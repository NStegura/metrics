package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/NStegura/metrics/internal/monitoring/pprof"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
)

const (
	timeoutShutdown = time.Second * 10
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func configureLogger(config *metricsapi.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
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
		return fmt.Errorf("failed to parse config: %w", err)
	}
	logger, err := configureLogger(config)
	if err != nil {
		return err
	}

	db, err := repo.New(
		ctx,
		config.DatabaseDSN,
		config.StoreInterval,
		config.FileStoragePath,
		config.Restore,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer logger.Info("closed DB")
		defer wg.Done()
		<-ctx.Done()

		db.Shutdown(ctx)
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

	go func(errs chan<- error) {
		if err = pprof.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("listen and serve pprof has failed: %w", err)
		}
	}(componentsErrs)

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
	if err := runRest(); err != nil {
		log.Fatal(err)
	}
}
