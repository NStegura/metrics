package agent

import (
	"encoding/json"
	"github.com/NStegura/metrics/internal/app/agent/models"
	"github.com/NStegura/metrics/internal/clients/metric"
	"github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	Alloc         models.MetricName = "Alloc"
	BuckHashSys   models.MetricName = "BuckHashSys"
	Frees         models.MetricName = "Frees"
	GCCPUFraction models.MetricName = "GCCPUFraction"
	GCSys         models.MetricName = "GCSys"
	HeapAlloc     models.MetricName = "HeapAlloc"
	HeapIdle      models.MetricName = "HeapIdle"
	HeapInuse     models.MetricName = "HeapInuse"
	HeapObjects   models.MetricName = "HeapObjects"
	HeapReleased  models.MetricName = "HeapReleased"
	HeapSys       models.MetricName = "HeapSys"
	LastGC        models.MetricName = "LastGC"
	Lookups       models.MetricName = "Lookups"
	MCacheInuse   models.MetricName = "MCacheInuse"
	MCacheSys     models.MetricName = "MCacheSys"
	MSpanInuse    models.MetricName = "MSpanInuse"
	MSpanSys      models.MetricName = "MSpanSys"
	Mallocs       models.MetricName = "Mallocs"
	NextGC        models.MetricName = "NextGC"
	NumGC         models.MetricName = "NumGC"
	NumForcedGC   models.MetricName = "NumForcedGC"
	OtherSys      models.MetricName = "OtherSys"
	PauseTotalNs  models.MetricName = "PauseTotalNs"
	StackInuse    models.MetricName = "StackInuse"
	StackSys      models.MetricName = "StackSys"
	Sys           models.MetricName = "Sys"
	TotalAlloc    models.MetricName = "TotalAlloc"

	RandomValue models.MetricName = "RandomValue"
	PollCount   models.MetricName = "PollCount"

	Gauge   models.MetricType = "gauge"
	Counter models.MetricType = "counter"
)

type Agent struct {
	config *Config

	logger *logrus.Logger
}

func New(config *Config, logger *logrus.Logger) *Agent {
	return &Agent{
		config: config,
		logger: logger,
	}
}

func (ag *Agent) Start() error {
	metricsCli := metric.New(ag.config.HTTPAddr)

	var mu sync.Mutex
	var metrics models.Metrics

	go func() {
		var counter int64 = 0
		for {
			counter++
			time.Sleep(ag.config.PollInterval)

			stats := runtime.MemStats{}
			runtime.ReadMemStats(&stats)
			mu.Lock()
			metrics = getMetricsFromStats(stats, counter)
			mu.Unlock()
		}
	}()

	for {
		time.Sleep(ag.config.ReportInterval)
		mu.Lock()
		for _, m := range metrics.GaugeMetrics {
			ag.logger.Info(*m)
			jsonBody, err := json.Marshal(m)
			if err != nil {
				ag.logger.Fatal(err)
			}
			err = metricsCli.UpdateMetric(jsonBody)
			if err != nil {
				ag.logger.Error(err)
			}
		}
		for _, m := range metrics.CounterMetrics {
			ag.logger.Info(*m)
			jsonBody, err := json.Marshal(m)
			if err != nil {
				ag.logger.Fatal(err)
			}
			err = metricsCli.UpdateMetric(jsonBody)
			if err != nil {
				ag.logger.Error(err)
			}
		}
		mu.Unlock()
	}
}

func getMetricsFromStats(stats runtime.MemStats, counter int64) models.Metrics {
	gaugeMetrics := make(map[models.MetricName]*models.GaugeMetric, 27)
	counterMetrics := make(map[models.MetricName]*models.CounterMetric, 1)
	metrics := models.Metrics{GaugeMetrics: gaugeMetrics, CounterMetrics: counterMetrics}

	metrics.GaugeMetrics[Alloc] = &models.GaugeMetric{Name: Alloc, Type: Gauge, Value: float64(stats.Alloc)}
	metrics.GaugeMetrics[BuckHashSys] = &models.GaugeMetric{Name: BuckHashSys, Type: Gauge, Value: float64(stats.BuckHashSys)}
	metrics.GaugeMetrics[Frees] = &models.GaugeMetric{Name: Frees, Type: Gauge, Value: float64(stats.Frees)}
	metrics.GaugeMetrics[GCCPUFraction] = &models.GaugeMetric{Name: GCCPUFraction, Type: Gauge, Value: float64(stats.GCCPUFraction)}
	metrics.GaugeMetrics[GCSys] = &models.GaugeMetric{Name: GCSys, Type: Gauge, Value: float64(stats.GCSys)}
	metrics.GaugeMetrics[HeapAlloc] = &models.GaugeMetric{Name: HeapAlloc, Type: Gauge, Value: float64(stats.HeapAlloc)}
	metrics.GaugeMetrics[HeapIdle] = &models.GaugeMetric{Name: HeapIdle, Type: Gauge, Value: float64(stats.HeapIdle)}
	metrics.GaugeMetrics[HeapInuse] = &models.GaugeMetric{Name: HeapInuse, Type: Gauge, Value: float64(stats.HeapInuse)}
	metrics.GaugeMetrics[HeapObjects] = &models.GaugeMetric{Name: HeapObjects, Type: Gauge, Value: float64(stats.HeapObjects)}
	metrics.GaugeMetrics[HeapReleased] = &models.GaugeMetric{Name: HeapReleased, Type: Gauge, Value: float64(stats.HeapReleased)}
	metrics.GaugeMetrics[HeapSys] = &models.GaugeMetric{Name: HeapSys, Type: Gauge, Value: float64(stats.HeapSys)}
	metrics.GaugeMetrics[LastGC] = &models.GaugeMetric{Name: LastGC, Type: Gauge, Value: float64(stats.LastGC)}
	metrics.GaugeMetrics[Lookups] = &models.GaugeMetric{Name: Lookups, Type: Gauge, Value: float64(stats.Lookups)}
	metrics.GaugeMetrics[MCacheInuse] = &models.GaugeMetric{Name: MCacheInuse, Type: Gauge, Value: float64(stats.MCacheInuse)}
	metrics.GaugeMetrics[MCacheSys] = &models.GaugeMetric{Name: MCacheSys, Type: Gauge, Value: float64(stats.MCacheSys)}
	metrics.GaugeMetrics[MSpanInuse] = &models.GaugeMetric{Name: MSpanInuse, Type: Gauge, Value: float64(stats.MSpanInuse)}
	metrics.GaugeMetrics[MSpanSys] = &models.GaugeMetric{Name: MSpanSys, Type: Gauge, Value: float64(stats.MSpanSys)}
	metrics.GaugeMetrics[Mallocs] = &models.GaugeMetric{Name: Mallocs, Type: Gauge, Value: float64(stats.Mallocs)}
	metrics.GaugeMetrics[NextGC] = &models.GaugeMetric{Name: NextGC, Type: Gauge, Value: float64(stats.NextGC)}
	metrics.GaugeMetrics[NumGC] = &models.GaugeMetric{Name: NumGC, Type: Gauge, Value: float64(stats.NumForcedGC)}
	metrics.GaugeMetrics[NumForcedGC] = &models.GaugeMetric{Name: NumForcedGC, Type: Gauge, Value: float64(stats.NumForcedGC)}
	metrics.GaugeMetrics[OtherSys] = &models.GaugeMetric{Name: OtherSys, Type: Gauge, Value: float64(stats.OtherSys)}
	metrics.GaugeMetrics[PauseTotalNs] = &models.GaugeMetric{Name: PauseTotalNs, Type: Gauge, Value: float64(stats.PauseTotalNs)}
	metrics.GaugeMetrics[StackInuse] = &models.GaugeMetric{Name: StackInuse, Type: Gauge, Value: float64(stats.StackInuse)}
	metrics.GaugeMetrics[StackSys] = &models.GaugeMetric{Name: StackSys, Type: Gauge, Value: float64(stats.StackSys)}
	metrics.GaugeMetrics[Sys] = &models.GaugeMetric{Name: Sys, Type: Gauge, Value: float64(stats.Sys)}
	metrics.GaugeMetrics[TotalAlloc] = &models.GaugeMetric{Name: TotalAlloc, Type: Gauge, Value: float64(stats.TotalAlloc)}

	metrics.GaugeMetrics[RandomValue] = &models.GaugeMetric{Name: RandomValue, Type: Gauge, Value: rand.Float64()}
	metrics.CounterMetrics[PollCount] = &models.CounterMetric{Name: PollCount, Type: Counter, Value: counter}
	return metrics
}
