package dal

import (
	"github.com/NStegura/metrics/internal/custom_errors"
	"github.com/NStegura/metrics/internal/dal/models"
	"github.com/sirupsen/logrus"
)

type repository struct {
	gaugeMetrics   map[string]*models.GaugeMetric
	counterMetrics map[string]*models.CounterMetric

	logger *logrus.Logger
}

func New() *repository {
	gaugeMetrics := map[string]*models.GaugeMetric{}
	counterMetrics := map[string]*models.CounterMetric{}
	return &repository{gaugeMetrics, counterMetrics, logrus.New()}
}

func (r *repository) GetCounterMetric(name string) (cm models.CounterMetric, err error) {
	metric, ok := r.counterMetrics[name]
	if ok != true {
		err = custom_errors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateCounterMetric(name string, mType string, value int64) {
	r.counterMetrics[name] = &models.CounterMetric{Name: name, Type: mType, Value: value}
}

func (r *repository) UpdateCounterMetric(name string, value int64) error {
	metric, ok := r.counterMetrics[name]
	if ok != true {
		return custom_errors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *repository) GetGaugeMetric(name string) (cm models.GaugeMetric, err error) {
	metric, ok := r.gaugeMetrics[name]
	if ok != true {
		err = custom_errors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateGaugeMetric(name string, mType string, value float64) {
	r.gaugeMetrics[name] = &models.GaugeMetric{Name: name, Type: mType, Value: value}
}

func (r *repository) UpdateGaugeMetric(name string, value float64) error {
	metric, ok := r.gaugeMetrics[name]
	if ok != true {
		return custom_errors.ErrNotFound
	}
	metric.Value = value
	return nil
}

func (r *repository) LogRepo() {
	r.logger.Info("gaugeMetrics: ---------------------------------")
	for name, metric := range r.gaugeMetrics {
		r.logger.Infof("name: %s, metric: %v", name, *metric)
	}
	r.logger.Info("counterMetrics: ---------------------------------")
	for name, metric := range r.counterMetrics {
		r.logger.Infof("name: %s metric: %v", name, *metric)
	}
}
