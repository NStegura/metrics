package business

import (
	"errors"
	blModels "github.com/NStegura/metrics/internal/business/models"
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/sirupsen/logrus"
	"sort"
)

type bll struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository) *bll {
	return &bll{repo: repo, logger: logrus.New()}
}

func (bll *bll) GetGaugeMetric(mName string) (float64, error) {
	gm, err := bll.repo.GetGaugeMetric(mName)
	if errors.Is(err, customerrors.ErrNotFound) {
		bll.logger.Warning(err) //debug
		return 0, err
	}

	return gm.Value, nil
}

func (bll *bll) UpdateGaugeMetric(gmReq blModels.GaugeMetric) (err error) {
	_, err = bll.repo.GetGaugeMetric(gmReq.Name)
	if errors.Is(err, customerrors.ErrNotFound) {
		bll.repo.CreateGaugeMetric(gmReq.Name, gmReq.Type, gmReq.Value)
	}
	err = bll.repo.UpdateGaugeMetric(gmReq.Name, gmReq.Value)

	// debug
	bll.repo.LogRepo()
	return
}

func (bll *bll) GetCounterMetric(mName string) (int64, error) {
	cm, err := bll.repo.GetCounterMetric(mName)
	if errors.Is(err, customerrors.ErrNotFound) {
		bll.logger.Warning(err) //debug
		return 0, err
	}

	return cm.Value, nil
}

func (bll *bll) UpdateCounterMetric(cmReq blModels.CounterMetric) (err error) {
	cm, err := bll.repo.GetCounterMetric(cmReq.Name)
	if errors.Is(err, customerrors.ErrNotFound) {
		bll.repo.CreateCounterMetric(cmReq.Name, cmReq.Type, cmReq.Value)
	}
	newVal := cm.Value + cmReq.Value
	err = bll.repo.UpdateCounterMetric(cmReq.Name, newVal)

	// debug
	bll.repo.LogRepo()
	return
}

func (bll *bll) GetAllMetrics() ([]blModels.GaugeMetric, []blModels.CounterMetric) {
	gaugeMetrics := make([]blModels.GaugeMetric, 0, 26)
	counterMetrics := make([]blModels.CounterMetric, 0, 1)

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
