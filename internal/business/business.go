package business

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"

	blModels "github.com/NStegura/metrics/internal/business/models"
	"github.com/NStegura/metrics/internal/customerrors"
)

const (
	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
)

// bll бизнес слой.
type bll struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *bll {
	return &bll{repo: repo, logger: logger}
}

// GetGaugeMetric получает gauge метрику по имени.
func (bll *bll) GetGaugeMetric(ctx context.Context, mName string) (float64, error) {
	gm, err := bll.repo.GetGaugeMetric(ctx, mName)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return 0, fmt.Errorf("gauge metric not found: %w", err)
		}
		return 0, fmt.Errorf("failed to get gauge metric, %w", err)
	}

	return gm.Value, nil
}

// UpdateGaugeMetric обновляет gauge метрику.
func (bll *bll) UpdateGaugeMetric(ctx context.Context, gmReq blModels.GaugeMetric) (err error) {
	_, err = bll.repo.GetGaugeMetric(ctx, gmReq.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			err := bll.repo.CreateGaugeMetric(ctx, gmReq.Name, gmReq.Type, gmReq.Value)
			if err != nil {
				return fmt.Errorf("create gauge metric failed, %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get gauge metric, %w", err)
	}
	err = bll.repo.UpdateGaugeMetric(ctx, gmReq.Name, gmReq.Value)
	return
}

// GetCounterMetric получает counter метрику по имени.
func (bll *bll) GetCounterMetric(ctx context.Context, mName string) (int64, error) {
	cm, err := bll.repo.GetCounterMetric(ctx, mName)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return 0, fmt.Errorf("counter metric not found: %w", err)
		}
		return 0, fmt.Errorf("failed to get counter metric, %w", err)
	}

	return cm.Value, nil
}

// UpdateCounterMetric обновляет counter метрику.
func (bll *bll) UpdateCounterMetric(ctx context.Context, cmReq blModels.CounterMetric) (err error) {
	cm, err := bll.repo.GetCounterMetric(ctx, cmReq.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			err := bll.repo.CreateCounterMetric(ctx, cmReq.Name, cmReq.Type, cmReq.Value)
			if err != nil {
				return fmt.Errorf("create counter metric failed, %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get counter metric, %w", err)
	}

	newVal := cm.Value + cmReq.Value
	err = bll.repo.UpdateCounterMetric(ctx, cmReq.Name, newVal)
	return
}

// GetAllMetrics получает все метрики.
func (bll *bll) GetAllMetrics(ctx context.Context) ([]blModels.GaugeMetric, []blModels.CounterMetric, error) {
	gaugeMetrics := make([]blModels.GaugeMetric, 0, countGaugeMetrics)
	counterMetrics := make([]blModels.CounterMetric, 0, countCounterMetrics)

	gms, cms, err := bll.repo.GetAllMetrics(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all metrics, %w", err)
	}

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

	return gaugeMetrics, counterMetrics, nil
}

// Ping проверяет работу сервера.
func (bll *bll) Ping(ctx context.Context) error {
	err := bll.repo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping repo %w", err)
	}
	return nil
}
