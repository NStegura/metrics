package repo

import (
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
	"time"
)

type Metrics struct {
	GaugeMetrics   map[string]*models.GaugeMetric   `json:"gauge_metrics"`
	CounterMetrics map[string]*models.CounterMetric `json:"counter_metrics"`
}

type repository struct {
	m               *Metrics
	storeInterval   time.Duration
	fileStoragePath string
	restore         bool
	synchronously   bool

	logger *logrus.Logger
}

func New(
	storeInterval time.Duration,
	fileStoragePath string,
	restore bool,
	logger *logrus.Logger,
) *repository {

	metrics := Metrics{
		map[string]*models.GaugeMetric{},
		map[string]*models.CounterMetric{},
	}

	if restore {
		logger.Info("LoadBackup")
		metrics, err := LoadBackup(fileStoragePath)
		if err != nil {
			logger.Warningf("backup load err, %s", err)
		}
		return &repository{
			&metrics,
			storeInterval,
			fileStoragePath,
			restore,
			storeInterval == 0,
			logger,
		}
	}
	return &repository{
		&metrics,
		storeInterval,
		fileStoragePath,
		restore,
		storeInterval == 0,
		logger,
	}
}

func (r *repository) GetCounterMetric(name string) (cm models.CounterMetric, err error) {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateCounterMetric(name string, mType string, value int64) {
	r.m.CounterMetrics[name] = &models.CounterMetric{Name: name, Type: mType, Value: value}
}

func (r *repository) UpdateCounterMetric(name string, value int64) error {
	metric, ok := r.m.CounterMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value

	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}

	return nil
}

func (r *repository) GetGaugeMetric(name string) (cm models.GaugeMetric, err error) {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		err = customerrors.ErrNotFound
		return
	}
	return *metric, nil
}

func (r *repository) CreateGaugeMetric(name string, mType string, value float64) {
	r.m.GaugeMetrics[name] = &models.GaugeMetric{Name: name, Type: mType, Value: value}

	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
}

func (r *repository) UpdateGaugeMetric(name string, value float64) error {
	metric, ok := r.m.GaugeMetrics[name]
	if !ok {
		return customerrors.ErrNotFound
	}
	metric.Value = value

	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
	return nil
}

func (r *repository) GetAllMetrics() ([]models.GaugeMetric, []models.CounterMetric) {
	gaugeMetrics := make([]models.GaugeMetric, 0, 26)
	counterMetrics := make([]models.CounterMetric, 0, 1)
	for _, gMetric := range r.m.GaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, *gMetric)
	}
	for _, cMetric := range r.m.CounterMetrics {
		counterMetrics = append(counterMetrics, *cMetric)
	}
	return gaugeMetrics, counterMetrics
}

func (r *repository) Shutdown() {
	r.logger.Info("Repo shutdown")
	err := r.MakeBackup()
	if err != nil {
		r.logger.Warningf("Last backup save err, %s", err)
	}
}
