package agent

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/NStegura/metrics/internal/app/agent/models"
	"github.com/NStegura/metrics/internal/clients/metric"
	"github.com/sirupsen/logrus"
)

const (
	alloc         models.MetricName = "Alloc"
	buckHashSys   models.MetricName = "BuckHashSys"
	frees         models.MetricName = "Frees"
	gccpuFraction models.MetricName = "GCCPUFraction"
	gcSys         models.MetricName = "GCSys"
	heapAlloc     models.MetricName = "HeapAlloc"
	heapIdle      models.MetricName = "HeapIdle"
	heapInuse     models.MetricName = "HeapInuse"
	heapObjects   models.MetricName = "HeapObjects"
	heapReleased  models.MetricName = "HeapReleased"
	heapSys       models.MetricName = "HeapSys"
	lastGC        models.MetricName = "LastGC"
	lookups       models.MetricName = "Lookups"
	mCacheInuse   models.MetricName = "MCacheInuse"
	mCacheSys     models.MetricName = "MCacheSys"
	mSpanInuse    models.MetricName = "MSpanInuse"
	mSpanSys      models.MetricName = "MSpanSys"
	mallocs       models.MetricName = "Mallocs"
	nextGC        models.MetricName = "NextGC"
	numGC         models.MetricName = "NumGC"
	numForcedGC   models.MetricName = "NumForcedGC"
	otherSys      models.MetricName = "OtherSys"
	pauseTotalNs  models.MetricName = "PauseTotalNs"
	stackInuse    models.MetricName = "StackInuse"
	stackSys      models.MetricName = "StackSys"
	sys           models.MetricName = "Sys"
	totalAlloc    models.MetricName = "TotalAlloc"

	randomValue models.MetricName = "RandomValue"
	pollCount   models.MetricName = "PollCount"

	gauge    models.MetricType = "gauge"
	counterT models.MetricType = "counter"

	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
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
				ag.logger.Error(err)
			}
			err = metricsCli.UpdateMetric(jsonBody, "gzip")
			if err != nil {
				ag.logger.Error(err)
			}
		}
		for _, m := range metrics.CounterMetrics {
			ag.logger.Info(*m)
			jsonBody, err := json.Marshal(m)
			if err != nil {
				ag.logger.Error(err)
			}
			err = metricsCli.UpdateMetric(jsonBody, "gzip")
			if err != nil {
				ag.logger.Error(err)
			}
		}
		mu.Unlock()
	}
}

func getMetricsFromStats(stats runtime.MemStats, counter int64) models.Metrics {
	gaugeMetrics := make(map[models.MetricName]*models.GaugeMetric, countGaugeMetrics)
	counterMetrics := make(map[models.MetricName]*models.CounterMetric, countCounterMetrics)
	metrics := models.Metrics{GaugeMetrics: gaugeMetrics, CounterMetrics: counterMetrics}

	metrics.GaugeMetrics[alloc] = &models.GaugeMetric{Name: alloc, Type: gauge, Value: float64(stats.Alloc)}
	metrics.GaugeMetrics[buckHashSys] = &models.GaugeMetric{
		Name: buckHashSys, Type: gauge, Value: float64(stats.BuckHashSys)}
	metrics.GaugeMetrics[frees] = &models.GaugeMetric{Name: frees, Type: gauge, Value: float64(stats.Frees)}
	metrics.GaugeMetrics[gccpuFraction] = &models.GaugeMetric{
		Name: gccpuFraction, Type: gauge, Value: stats.GCCPUFraction}
	metrics.GaugeMetrics[gcSys] = &models.GaugeMetric{Name: gcSys, Type: gauge, Value: float64(stats.GCSys)}
	metrics.GaugeMetrics[heapAlloc] = &models.GaugeMetric{Name: heapAlloc, Type: gauge, Value: float64(stats.HeapAlloc)}
	metrics.GaugeMetrics[heapIdle] = &models.GaugeMetric{Name: heapIdle, Type: gauge, Value: float64(stats.HeapIdle)}
	metrics.GaugeMetrics[heapInuse] = &models.GaugeMetric{Name: heapInuse, Type: gauge, Value: float64(stats.HeapInuse)}
	metrics.GaugeMetrics[heapObjects] = &models.GaugeMetric{
		Name: heapObjects, Type: gauge, Value: float64(stats.HeapObjects)}
	metrics.GaugeMetrics[heapReleased] = &models.GaugeMetric{
		Name: heapReleased, Type: gauge, Value: float64(stats.HeapReleased)}
	metrics.GaugeMetrics[heapSys] = &models.GaugeMetric{Name: heapSys, Type: gauge, Value: float64(stats.HeapSys)}
	metrics.GaugeMetrics[lastGC] = &models.GaugeMetric{Name: lastGC, Type: gauge, Value: float64(stats.LastGC)}
	metrics.GaugeMetrics[lookups] = &models.GaugeMetric{Name: lookups, Type: gauge, Value: float64(stats.Lookups)}
	metrics.GaugeMetrics[mCacheInuse] = &models.GaugeMetric{
		Name: mCacheInuse, Type: gauge, Value: float64(stats.MCacheInuse)}
	metrics.GaugeMetrics[mCacheSys] = &models.GaugeMetric{Name: mCacheSys, Type: gauge, Value: float64(stats.MCacheSys)}
	metrics.GaugeMetrics[mSpanInuse] = &models.GaugeMetric{
		Name: mSpanInuse, Type: gauge, Value: float64(stats.MSpanInuse)}
	metrics.GaugeMetrics[mSpanSys] = &models.GaugeMetric{Name: mSpanSys, Type: gauge, Value: float64(stats.MSpanSys)}
	metrics.GaugeMetrics[mallocs] = &models.GaugeMetric{Name: mallocs, Type: gauge, Value: float64(stats.Mallocs)}
	metrics.GaugeMetrics[nextGC] = &models.GaugeMetric{Name: nextGC, Type: gauge, Value: float64(stats.NextGC)}
	metrics.GaugeMetrics[numGC] = &models.GaugeMetric{Name: numGC, Type: gauge, Value: float64(stats.NumForcedGC)}
	metrics.GaugeMetrics[numForcedGC] = &models.GaugeMetric{
		Name: numForcedGC, Type: gauge, Value: float64(stats.NumForcedGC)}
	metrics.GaugeMetrics[otherSys] = &models.GaugeMetric{Name: otherSys, Type: gauge, Value: float64(stats.OtherSys)}
	metrics.GaugeMetrics[pauseTotalNs] = &models.GaugeMetric{
		Name: pauseTotalNs, Type: gauge, Value: float64(stats.PauseTotalNs)}
	metrics.GaugeMetrics[stackInuse] = &models.GaugeMetric{
		Name: stackInuse, Type: gauge, Value: float64(stats.StackInuse)}
	metrics.GaugeMetrics[stackSys] = &models.GaugeMetric{Name: stackSys, Type: gauge, Value: float64(stats.StackSys)}
	metrics.GaugeMetrics[sys] = &models.GaugeMetric{Name: sys, Type: gauge, Value: float64(stats.Sys)}
	metrics.GaugeMetrics[totalAlloc] = &models.GaugeMetric{
		Name: totalAlloc, Type: gauge, Value: float64(stats.TotalAlloc)}

	metrics.GaugeMetrics[randomValue] = &models.GaugeMetric{Name: randomValue, Type: gauge, Value: rand.Float64()}
	metrics.CounterMetrics[pollCount] = &models.CounterMetric{Name: pollCount, Type: counterT, Value: counter}
	return metrics
}
