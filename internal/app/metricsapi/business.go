package metricsapi

import (
	"context"

	blModels "github.com/NStegura/metrics/internal/business/models"
)

type Bll interface {
	GetGaugeMetric(context.Context, string) (float64, error)
	UpdateGaugeMetric(context.Context, blModels.GaugeMetric) error
	GetCounterMetric(context.Context, string) (int64, error)
	UpdateCounterMetric(context.Context, blModels.CounterMetric) error
	GetAllMetrics(context.Context) ([]blModels.GaugeMetric, []blModels.CounterMetric, error)

	Ping(context.Context) error
}
