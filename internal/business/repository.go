package business

import (
	"context"

	"github.com/NStegura/metrics/internal/repo/models"
)

type Repository interface {
	GetCounterMetric(ctx context.Context, name string) (models.CounterMetric, error)
	CreateCounterMetric(ctx context.Context, name string, mType string, value int64)
	UpdateCounterMetric(ctx context.Context, name string, value int64) error
	GetGaugeMetric(context.Context, string) (models.GaugeMetric, error)
	CreateGaugeMetric(ctx context.Context, name string, mType string, value float64)
	UpdateGaugeMetric(ctx context.Context, name string, value float64) error
	GetAllMetrics(ctx context.Context) ([]models.GaugeMetric, []models.CounterMetric)

	Ping(ctx context.Context) error
}
