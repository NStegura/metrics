package metricsapi

import blModels "github.com/NStegura/metrics/internal/bll/models"

type Bll interface {
	GetGaugeMetric(string) (float64, error)
	UpdateGaugeMetric(blModels.GaugeMetric) error
	GetCounterMetric(string) (int64, error)
	UpdateCounterMetric(blModels.CounterMetric) error
	GetAllMetrics() ([]blModels.GaugeMetric, []blModels.CounterMetric)
}
