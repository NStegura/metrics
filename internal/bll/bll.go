package bll

import (
	"errors"
	blModels "github.com/NStegura/metrics/internal/bll/models"
	"github.com/NStegura/metrics/internal/custom_errors"
	"github.com/sirupsen/logrus"
)

type bll struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository) *bll {
	return &bll{repo: repo}
}

func (bll *bll) UpdateGaugeMetric(gmReq blModels.GaugeMetric) (err error) {
	_, err = bll.repo.GetGaugeMetric(gmReq.Name)
	if errors.Is(err, custom_errors.ErrNotFound) {
		bll.repo.CreateGaugeMetric(gmReq.Name, gmReq.Type, gmReq.Value)
	}
	err = bll.repo.UpdateGaugeMetric(gmReq.Name, gmReq.Value)

	// debug
	bll.repo.LogRepo()
	return
}

func (bll *bll) UpdateCounterMetric(cmReq blModels.CounterMetric) (err error) {
	cm, err := bll.repo.GetCounterMetric(cmReq.Name)
	if errors.Is(err, custom_errors.ErrNotFound) {
		bll.repo.CreateCounterMetric(cmReq.Name, cmReq.Type, cmReq.Value)
	}
	newVal := cm.Value + cmReq.Value
	err = bll.repo.UpdateCounterMetric(cmReq.Name, newVal)

	// debug
	bll.repo.LogRepo()
	return
}
