package mem

import (
	"context"

	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
)

const (
	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
)

type Metrics struct {
	GaugeMetrics   map[string]*models.GaugeMetric   `json:"gauge_metrics"`
	CounterMetrics map[string]*models.CounterMetric `json:"counter_metrics"`
}

type InMemoryRepo struct {
	m *Metrics

	logger *logrus.Logger
}

func NewInMemoryRepo(logger *logrus.Logger) (*InMemoryRepo, error) {
	return &InMemoryRepo{nil, logger}, nil
}

func (r *InMemoryRepo) GetCounterMetric(ctx context.Context, name string) (cm models.CounterMetric, err error) {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *InMemoryRepo) CreateCounterMetric(ctx context.Context, name string, mType string, value int64) {
	r.m.CounterMetrics[name] = &models.CounterMetric{Name: name, Type: mType, Value: value}
}

func (r *InMemoryRepo) UpdateCounterMetric(ctx context.Context, name string, value int64) error {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *InMemoryRepo) GetGaugeMetric(ctx context.Context, name string) (cm models.GaugeMetric, err error) {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *InMemoryRepo) CreateGaugeMetric(ctx context.Context, name string, mType string, value float64) {
	r.m.GaugeMetrics[name] = &models.GaugeMetric{Name: name, Type: mType, Value: value}
}

func (r *InMemoryRepo) UpdateGaugeMetric(ctx context.Context, name string, value float64) error {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *InMemoryRepo) GetAllMetrics(ctx context.Context) ([]models.GaugeMetric, []models.CounterMetric) {
	gaugeMetrics := make([]models.GaugeMetric, 0, countGaugeMetrics)
	counterMetrics := make([]models.CounterMetric, 0, countCounterMetrics)
	for _, gMetric := range r.m.GaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, *gMetric)
	}
	for _, cMetric := range r.m.CounterMetrics {
		counterMetrics = append(counterMetrics, *cMetric)
	}
	return gaugeMetrics, counterMetrics
}

func (r *InMemoryRepo) Init(ctx context.Context) error {
	r.logger.Info("Init repo")
	r.m = &Metrics{
		map[string]*models.GaugeMetric{},
		map[string]*models.CounterMetric{},
	}
	return nil
}

func (r *InMemoryRepo) Shutdown(ctx context.Context) {
	r.logger.Info("Repo shutdown")
}

func (r *InMemoryRepo) Ping(ctx context.Context) error {
	r.logger.Info("Pong")
	return nil
}
