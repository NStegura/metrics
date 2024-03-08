package metricsapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/metricsapi/models"
	blModels "github.com/NStegura/metrics/internal/business/models"
	"github.com/NStegura/metrics/internal/customerrors"
)

type contentType string
type URLParam string
type metricType string

const (
	contType string = "Content-Type"

	textHTML contentType = "text/html"

	mName  URLParam = "mName"
	mValue URLParam = "mValue"

	gauge   metricType = "gauge"
	counter metricType = "counter"

	timeout = 5 * time.Second
)

// APIServer хранит сущности для работы сервера.
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

// Start запускает сервер.
func (s *APIServer) Start() error {
	s.configRouter()

	s.logger.Infof("starting APIServer %s", s.config.BindAddr)
	if err := http.ListenAndServe(s.config.BindAddr, s.router); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *APIServer) configRouter() {
	s.router.Use(s.requestLogger)
	s.router.Use(s.gzipMiddleware)
	s.router.Use(s.hashValidation)

	s.router.Get(`/`, s.getAllMetrics())
	s.router.Post(`/updates/`, s.updateAllMetrics())

	s.router.Route(`/value`, func(r chi.Router) {
		r.Post(`/`, s.getMetric())
		r.Route(`/gauge`, func(r chi.Router) {
			r.Get(`/{mName}`, s.getGaugeMetric())
		})
		r.Route(`/counter`, func(r chi.Router) {
			r.Get(`/{mName}`, s.getCounterMetric())
		})
	})

	s.router.Route(`/update`, func(r chi.Router) {
		r.Post(`/`, s.updateMetric())
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

	s.router.Get(`/ping`, s.ping())
}

func (s *APIServer) getAllMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		gms, cms, err := s.bll.GetAllMetrics(ctx)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var sb strings.Builder

		for _, m := range gms {
			sb.WriteString(fmt.Sprintf("%s: %v\r\n", m.Name, m.Value))
		}
		for _, m := range cms {
			sb.WriteString(fmt.Sprintf("%s: %v\r\n", m.Name, m.Value))
		}
		w.Header().Set(contType, string(textHTML))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(sb.String())); err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *APIServer) updateAllMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var metrics models.MetricsList

		if err := easyjson.UnmarshalFromReader(r.Body, &metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, metric := range metrics {
			switch metric.MType {
			case string(gauge):
				if metric.Value == nil {
					http.Error(w, "gauge metric value null", http.StatusBadRequest)
					return
				}
				err := s.bll.UpdateGaugeMetric(
					ctx,
					blModels.GaugeMetric{
						Name: metric.ID, Type: metric.MType, Value: *metric.Value,
					})
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			case string(counter):
				if metric.Delta == nil {
					http.Error(w, "counter metric value null", http.StatusBadRequest)
					return
				}
				err := s.bll.UpdateCounterMetric(
					ctx,
					blModels.CounterMetric{
						Name: metric.ID, Type: metric.MType, Value: *metric.Delta,
					})
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			default:
				http.Error(w, "unknown metric type", http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) getCounterMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		mn := chi.URLParam(r, string(mName))
		if mn == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric, err := s.bll.GetCounterMetric(ctx, mn)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set(contType, string(textHTML))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(strconv.FormatInt(metric, 10))); err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *APIServer) updateCounterMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		mn := chi.URLParam(r, string(mName))
		mv := chi.URLParam(r, string(mValue))

		m, err := parseMetric(r.URL.Path, string(counter), mn, mv)

		if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cm, err := models.CastToCounter(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.bll.UpdateCounterMetric(ctx, blModels.CounterMetric(cm))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) getGaugeMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		mn := chi.URLParam(r, string(mName))
		if mn == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric, err := s.bll.GetGaugeMetric(ctx, mn)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set(contType, string(textHTML))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("%v", metric))); err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *APIServer) updateGaugeMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		mn := chi.URLParam(r, string(mName))
		mv := chi.URLParam(r, string(mValue))

		m, err := parseMetric(r.URL.Path, string(gauge), mn, mv)

		if errors.As(err, &customerrors.ParseURLError{URL: r.URL.Path}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		gm, err := models.CastToGauge(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.bll.UpdateGaugeMetric(ctx, blModels.GaugeMetric(gm))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) updateMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var metric models.Metrics

		if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case string(gauge):
			if metric.Value == nil {
				http.Error(w, "metric value null", http.StatusBadRequest)
				return
			}
			err := s.bll.UpdateGaugeMetric(
				ctx,
				blModels.GaugeMetric{
					Name: metric.ID, Type: metric.MType, Value: *metric.Value,
				})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		case string(counter):
			if metric.Delta == nil {
				http.Error(w, "metric value null", http.StatusBadRequest)
				return
			}
			err := s.bll.UpdateCounterMetric(
				ctx,
				blModels.CounterMetric{
					Name: metric.ID, Type: metric.MType, Value: *metric.Delta,
				})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		default:
			http.Error(w, "unknown metric type", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) getMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		var metric models.Metrics

		if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case string(gauge):
			gm, err := s.bll.GetGaugeMetric(ctx, metric.ID)
			if err != nil {
				if errors.Is(err, customerrors.ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			metric.Value = &gm
			s.writeJSONResp(metric, w)
		case string(counter):
			cm, err := s.bll.GetCounterMetric(ctx, metric.ID)
			if err != nil {
				if errors.Is(err, customerrors.ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			metric.Delta = &cm
			s.writeJSONResp(metric, w)
		default:
			http.Error(w, "unknown metric type", http.StatusBadRequest)
			return
		}
	}
}

func (s *APIServer) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		if err := s.bll.Ping(ctx); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) writeJSONResp(resp easyjson.Marshaler, w http.ResponseWriter) {
	w.Header().Set(contType, "application/json")

	jsonResp, err := easyjson.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		s.logger.Error(err)
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
