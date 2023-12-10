package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/NStegura/metrics/internal/repo/internal/db"
	"github.com/NStegura/metrics/internal/repo/internal/mem"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetCounterMetric(ctx context.Context, name string) (models.CounterMetric, error)
	CreateCounterMetric(ctx context.Context, name string, mType string, value int64)
	UpdateCounterMetric(ctx context.Context, name string, value int64) error
	GetGaugeMetric(ctx context.Context, name string) (models.GaugeMetric, error)
	CreateGaugeMetric(ctx context.Context, name string, mType string, value float64)
	UpdateGaugeMetric(ctx context.Context, name string, value float64) error
	GetAllMetrics(ctx context.Context) ([]models.GaugeMetric, []models.CounterMetric)

	Init(ctx context.Context) error
	Shutdown(ctx context.Context)
	Ping(ctx context.Context) error
}

func New(
	ctx context.Context,
	dbDSN string,
	storeInterval time.Duration,
	fileStoragePath string,
	restore bool,
	logger *logrus.Logger,
) (Repository, error) {
	if dbDSN != "" {
		logger.Info("Init db")
		repo, err := db.New(ctx, dbDSN, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create DB: %w", err)
		}
		return repo, nil
	}

	if restore {
		logger.Info("Init mem repo with backup")
		repo, err := mem.NewBackupRepo(storeInterval, fileStoragePath, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create mem repo with backup: %w", err)
		}
		return repo, nil
	}

	logger.Info("Init mem repo")
	repo, err := mem.NewInMemoryRepo(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create mem repo: %w", err)
	}
	return repo, nil
}
