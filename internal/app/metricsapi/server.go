package metricsapi

import (
	"errors"
	"github.com/NStegura/metrics/internal/app/metricsapi/models"
	blModels "github.com/NStegura/metrics/internal/bll/models"
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type APIServer struct {
	config *Config
	bll    Bll
	router *http.ServeMux

	logger *logrus.Logger
}

func New(config *Config, bll Bll) *APIServer {
	return &APIServer{
		config: config,
		bll:    bll,
		router: http.NewServeMux(),
		logger: logrus.New(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configRouter()
	s.logger.Info("starting APIServer")
	return http.ListenAndServe(`:8080`, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configRouter() {
	s.router.HandleFunc(`/update/gauge/`, s.updateGaugeMetric())
	s.router.HandleFunc(`/update/counter/`, s.updateCounterMetric())
	s.router.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
}

func (s *APIServer) updateCounterMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			m, err := parseMetric(r.URL.Path)
			if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if m.Value == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			cm, err := models.CastToCounter(m)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			s.logger.Info(cm) // debug
			err = s.bll.UpdateCounterMetric(blModels.CounterMetric(cm))
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			return
		}
	}
}

func (s *APIServer) updateGaugeMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			m, err := parseMetric(r.URL.Path)
			if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if m.Value == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			gm, err := models.CastToGauge(m)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			s.logger.Info(gm) // debug
			err = s.bll.UpdateGaugeMetric(blModels.GaugeMetric(gm))
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			return
		}
	}
}

func parseMetric(url string) (metric models.Metric, err error) {
	fullURLFragments := strings.Split(url, "/")

	if len(fullURLFragments) != 5 {
		err = &customerrors.ParseURLError{URL: url}
		return
	}

	values := fullURLFragments[2:]

	return models.Metric{
		Name:  values[1],
		Type:  values[0],
		Value: values[2],
	}, nil
}
