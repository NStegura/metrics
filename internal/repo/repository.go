package repo

import (
	"fmt"
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
	"strings"
)

type repository struct {
	gaugeMetrics   map[string]*models.GaugeMetric
	counterMetrics map[string]*models.CounterMetric

	logger *logrus.Logger
}

func New(logger *logrus.Logger) *repository {
	gaugeMetrics := map[string]*models.GaugeMetric{}
	counterMetrics := map[string]*models.CounterMetric{}
	return &repository{gaugeMetrics, counterMetrics, logger}
}

func (r *repository) GetCounterMetric(name string) (cm models.CounterMetric, err error) {
	metric, ok := r.counterMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateCounterMetric(name string, mType string, value int64) {
	r.counterMetrics[name] = &models.CounterMetric{Name: name, Type: mType, Value: value}
}

func (r *repository) UpdateCounterMetric(name string, value int64) error {
	metric, ok := r.counterMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *repository) GetGaugeMetric(name string) (cm models.GaugeMetric, err error) {
	metric, ok := r.gaugeMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateGaugeMetric(name string, mType string, value float64) {
	r.gaugeMetrics[name] = &models.GaugeMetric{Name: name, Type: mType, Value: value}
}

func (r *repository) UpdateGaugeMetric(name string, value float64) error {
	metric, ok := r.gaugeMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *repository) GetAllMetrics() ([]models.GaugeMetric, []models.CounterMetric) {
	gaugeMetrics := make([]models.GaugeMetric, 0, 26)
	counterMetrics := make([]models.CounterMetric, 0, 1)
	for _, gMetric := range r.gaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, *gMetric)
	}
	for _, cMetric := range r.counterMetrics {
		counterMetrics = append(counterMetrics, *cMetric)
	}
	return gaugeMetrics, counterMetrics
}

func (r *repository) LogRepo() {
	var sb strings.Builder

	sb.WriteString("\ngaugeMetrics: ---------------------------------\n\t")
	for name, metric := range r.gaugeMetrics {
		sb.WriteString(fmt.Sprintf("name: %s, metric: %v\n\t", name, *metric))
	}
	sb.WriteString("\ncounterMetrics: ---------------------------------\n\t")
	for name, metric := range r.counterMetrics {
		sb.WriteString(fmt.Sprintf("name: %s, metric: %v\n\t", name, *metric))
	}

	r.logger.Info(sb.String())
}
