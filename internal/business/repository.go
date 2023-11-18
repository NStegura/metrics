package business

import "github.com/NStegura/metrics/internal/repo/models"

type Repository interface {
	GetCounterMetric(name string) (models.CounterMetric, error)
	CreateCounterMetric(name string, mType string, value int64)
	UpdateCounterMetric(name string, value int64) error
	GetGaugeMetric(name string) (models.GaugeMetric, error)
	CreateGaugeMetric(name string, mType string, value float64)
	UpdateGaugeMetric(name string, value float64) error
	GetAllMetrics() ([]models.GaugeMetric, []models.CounterMetric)
}
