package repo

import (
	"time"

	"github.com/NStegura/metrics/internal/repo/internal/mem"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetCounterMetric(name string) (models.CounterMetric, error)
	CreateCounterMetric(name string, mType string, value int64)
	UpdateCounterMetric(name string, value int64) error
	GetGaugeMetric(name string) (models.GaugeMetric, error)
	CreateGaugeMetric(name string, mType string, value float64)
	UpdateGaugeMetric(name string, value float64) error
	GetAllMetrics() ([]models.GaugeMetric, []models.CounterMetric)
	Init() error
	Shutdown()
}

func New(
	storeInterval time.Duration,
	fileStoragePath string,
	restore bool,
	logger *logrus.Logger,
) Repository {
	if restore {
		logger.Info("Init mem repo with backup")
		return mem.NewBackupRepo(storeInterval, fileStoragePath, logger)
	}
	logger.Info("Init mem repo")
	return mem.NewInMemoryRepo(logger)
}
