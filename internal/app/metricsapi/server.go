package metricsapi

import (
	"errors"
	"fmt"
	"github.com/NStegura/metrics/internal/app/metricsapi/models"
	blModels "github.com/NStegura/metrics/internal/business/models"
	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type APIServer struct {
	config *Config
	bll    Bll
	router *chi.Mux

	logger *logrus.Logger
}

func New(config *Config, bll Bll, logger *logrus.Logger) *APIServer {
	return &APIServer{
		config: config,
		bll:    bll,
		router: chi.NewRouter(),
		logger: logger,
	}
}

func (s *APIServer) Start() error {
	s.configRouter()
	s.logger.Info("starting APIServer")
	return http.ListenAndServe(s.config.BindAddr, s.requestLogger(s.router))
}

func (s *APIServer) configRouter() {
	s.router.Get(`/`, s.getAllMetrics())

	s.router.Route(`/value`, func(r chi.Router) {
		r.Route(`/gauge`, func(r chi.Router) {
			r.Get(`/{mName}`, s.getGaugeMetric())
		})
		r.Route(`/counter`, func(r chi.Router) {
			r.Get(`/{mName}`, s.getCounterMetric())
		})
	})

	s.router.Route(`/update`, func(r chi.Router) {
		r.Post(`/{mType}/{mName}/{mValue}`, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
		r.Route(`/gauge`, func(r chi.Router) {
			r.Post(`/{mName}/{mValue}`, s.updateGaugeMetric())
		})
		r.Route(`/counter`, func(r chi.Router) {
			r.Post(`/{mName}/{mValue}`, s.updateCounterMetric())
		})
	})
}

func (s *APIServer) getAllMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sb strings.Builder
		gms, cms := s.bll.GetAllMetrics()

		for _, m := range gms {
			sb.WriteString(fmt.Sprintf("%s: %v\r\n", m.Name, m.Value))
		}
		for _, m := range cms {
			sb.WriteString(fmt.Sprintf("%s: %v\r\n", m.Name, m.Value))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sb.String()))
	}
}

func (s *APIServer) getCounterMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mName := chi.URLParam(r, "mName")
		if mName == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric, err := s.bll.GetCounterMetric(mName)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatInt(metric, 10)))
	}
}

func (s *APIServer) updateCounterMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")

		m, err := parseMetric(r.URL.Path, "counter", mName, mValue)

		if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cm, err := models.CastToCounter(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.bll.UpdateCounterMetric(blModels.CounterMetric(cm))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) getGaugeMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mName := chi.URLParam(r, "mName")
		if mName == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric, err := s.bll.GetGaugeMetric(mName)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", metric)))
	}
}

func (s *APIServer) updateGaugeMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")

		m, err := parseMetric(r.URL.Path, "gauge", mName, mValue)

		if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gm, err := models.CastToGauge(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.bll.UpdateGaugeMetric(blModels.GaugeMetric(gm))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func parseMetric(url string, mtype string, mName string, mValue string) (metric models.Metric, err error) {
	if mName == "" {
		err = &customerrors.ParseURLError{URL: url}
		return
	}

	return models.Metric{
		Name:  mName,
		Type:  mtype,
		Value: mValue,
	}, nil
}
