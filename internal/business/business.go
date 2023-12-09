package business

import (
	"errors"
	"fmt"
	"sort"

	blModels "github.com/NStegura/metrics/internal/business/models"
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/sirupsen/logrus"
)

const (
	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
)

type bll struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *bll {
	return &bll{repo: repo, logger: logger}
}

func (bll *bll) GetGaugeMetric(mName string) (float64, error) {
	gm, err := bll.repo.GetGaugeMetric(mName)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return 0, fmt.Errorf("gauge metric not found: %w", err)
		}
		return 0, fmt.Errorf("failed to get gauge metric, %w", err)
	}

	return gm.Value, nil
}

func (bll *bll) UpdateGaugeMetric(gmReq blModels.GaugeMetric) (err error) {
	_, err = bll.repo.GetGaugeMetric(gmReq.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			bll.repo.CreateGaugeMetric(gmReq.Name, gmReq.Type, gmReq.Value)
			return nil
		}
		return fmt.Errorf("failed to get gauge metric, %w", err)
	}
	err = bll.repo.UpdateGaugeMetric(gmReq.Name, gmReq.Value)
	return
}

func (bll *bll) GetCounterMetric(mName string) (int64, error) {
	cm, err := bll.repo.GetCounterMetric(mName)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return 0, fmt.Errorf("counter metric not found: %w", err)
		}
		return 0, fmt.Errorf("failed to get counter metric, %w", err)
	}

	return cm.Value, nil
}

func (bll *bll) UpdateCounterMetric(cmReq blModels.CounterMetric) (err error) {
	cm, err := bll.repo.GetCounterMetric(cmReq.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			bll.repo.CreateCounterMetric(cmReq.Name, cmReq.Type, cmReq.Value)
			return nil
		}
		return fmt.Errorf("failed to get counter metric, %w", err)
	}

	newVal := cm.Value + cmReq.Value
	err = bll.repo.UpdateCounterMetric(cmReq.Name, newVal)
	return
}

func (bll *bll) GetAllMetrics() ([]blModels.GaugeMetric, []blModels.CounterMetric) {
	gaugeMetrics := make([]blModels.GaugeMetric, 0, countGaugeMetrics)
	counterMetrics := make([]blModels.CounterMetric, 0, countCounterMetrics)

	gms, cms := bll.repo.GetAllMetrics()

	for _, gMetric := range gms {
		gaugeMetrics = append(gaugeMetrics, blModels.GaugeMetric{
			Name:  gMetric.Name,
			Type:  gMetric.Type,
			Value: gMetric.Value,
		})
	}
	for _, cMetric := range cms {
		counterMetrics = append(counterMetrics, blModels.CounterMetric{
			Name:  cMetric.Name,
			Type:  cMetric.Type,
			Value: cMetric.Value,
		})
	}
	if len(gaugeMetrics) > 1 {
		sort.Sort(blModels.ByName(gaugeMetrics))
	}

	return gaugeMetrics, counterMetrics
}
