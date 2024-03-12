package mem

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/NStegura/metrics/internal/repo/models"
)

const (
	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
)

type Metrics struct {
	GaugeMetrics   map[string]*models.GaugeMetric   `json:"gauge_metrics"`
	CounterMetrics map[string]*models.CounterMetric `json:"counter_metrics"`
}

// InMemoryRepo структура хранилища в памяти.
type InMemoryRepo struct {
	m *Metrics

	logger *logrus.Logger
}

func NewInMemoryRepo(logger *logrus.Logger) (*InMemoryRepo, error) {
	return &InMemoryRepo{
		&Metrics{
			map[string]*models.GaugeMetric{},
			map[string]*models.CounterMetric{},
		}, logger}, nil
}

// GetCounterMetric получает counter метрику по названию.
func (r *InMemoryRepo) GetCounterMetric(_ context.Context, name string) (cm models.CounterMetric, err error) {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

// CreateCounterMetric создает counter метрику.
func (r *InMemoryRepo) CreateCounterMetric(_ context.Context, name string, mType string, value int64) error {
	r.m.CounterMetrics[name] = &models.CounterMetric{Name: name, Type: mType, Value: value}
	return nil
}

// UpdateCounterMetric обновляет counter метрику.
func (r *InMemoryRepo) UpdateCounterMetric(_ context.Context, name string, value int64) error {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

// GetGaugeMetric получает gauge метрику.
func (r *InMemoryRepo) GetGaugeMetric(_ context.Context, name string) (cm models.GaugeMetric, err error) {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

// CreateGaugeMetric создает gauge метрику.
func (r *InMemoryRepo) CreateGaugeMetric(_ context.Context, name string, mType string, value float64) error {
	r.m.GaugeMetrics[name] = &models.GaugeMetric{Name: name, Type: mType, Value: value}
	return nil
}

// UpdateGaugeMetric обновляет gauge метрику.
func (r *InMemoryRepo) UpdateGaugeMetric(_ context.Context, name string, value float64) error {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

// GetAllMetrics получает все метрики.
func (r *InMemoryRepo) GetAllMetrics(_ context.Context) ([]models.GaugeMetric, []models.CounterMetric, error) {
	gaugeMetrics := make([]models.GaugeMetric, 0, countGaugeMetrics)
	counterMetrics := make([]models.CounterMetric, 0, countCounterMetrics)
	for _, gMetric := range r.m.GaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, *gMetric)
	}
	for _, cMetric := range r.m.CounterMetrics {
		counterMetrics = append(counterMetrics, *cMetric)
	}
	return gaugeMetrics, counterMetrics, nil
}

func (r *InMemoryRepo) Shutdown(_ context.Context) {
	r.logger.Info("Repo shutdown")
}

func (r *InMemoryRepo) Ping(_ context.Context) error {
	r.logger.Info("Pong")
	return nil
}
