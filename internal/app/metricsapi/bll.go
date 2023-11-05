package metricsapi

import blModels "github.com/NStegura/metrics/internal/bll/models"

type Bll interface {
	UpdateGaugeMetric(blModels.GaugeMetric) error
	UpdateCounterMetric(blModels.CounterMetric) error
}
