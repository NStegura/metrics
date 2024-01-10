package agent

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/metrics/internal/app/agent/models"
	"github.com/NStegura/metrics/internal/clients/metric"
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

	totalMemory     models.MetricName = "TotalMemory"
	freeMemory      models.MetricName = "FreeMemory"
	CPUutilization1 models.MetricName = "CPUutilization1"

	gauge    models.MetricType = "gauge"
	counterT models.MetricType = "counter"

	countGaugeMetrics   int = 27
	countCounterMetrics int = 1
	countGaugePsMetrics int = 3
)

type Agent struct {
	config     *Config
	metricsCli *metric.Client

	logger *logrus.Logger
}

func New(config *Config, metricsCli *metric.Client, logger *logrus.Logger) *Agent {
	return &Agent{
		config:     config,
		metricsCli: metricsCli,
		logger:     logger,
	}
}

func (ag *Agent) Start() error {
	var wg sync.WaitGroup

	metricsCh := ag.collectMetrics(&wg)

	ag.logger.Info(ag.config.RateLimit)
	for w := 1; w <= ag.config.RateLimit; w++ {
		ag.sendMetrics(w, &wg, metricsCh)
	}

	wg.Wait()
	return nil
}

func (ag *Agent) collectMetrics(wg *sync.WaitGroup) chan models.Metrics {
	metricsCh := make(chan models.Metrics, ag.config.RateLimit)
	pollTicker := time.NewTicker(ag.config.PollInterval)

	wg.Add(1)
	go func() {
		defer close(metricsCh)
		defer pollTicker.Stop()
		defer wg.Done()

		var counter int64 = 0

		for range pollTicker.C {
			counter++
			statMetrics := ag.getMetricsFromStats(counter)
			metricsCh <- statMetrics
			psMetrics, errs, maxErrs := ag.getPSMetrics()
			if len(errs) < maxErrs {
				metricsCh <- psMetrics
			}
		}
	}()

	return metricsCh
}

func (ag *Agent) sendMetrics(workerID int, wg *sync.WaitGroup, metricsCh <-chan models.Metrics) {
	ag.logger.Infof("start worker %v", workerID)
	reportTicker := time.NewTicker(ag.config.ReportInterval)

	wg.Add(1)
	go func() {
		defer reportTicker.Stop()
		defer wg.Done()

		for range reportTicker.C {
			err := ag.metricsCli.UpdateMetrics(metric.CastToMetrics(<-metricsCh), "gzip")
			if err != nil {
				ag.logger.Error(err)
			}
		}
	}()
}

func (ag *Agent) getMetricsFromStats(counter int64) models.Metrics {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

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

func (ag *Agent) getPSMetrics() (m models.Metrics, errs []error, maxErrs int) {
	maxErrsCount := 2

	gaugeMetrics := make(map[models.MetricName]*models.GaugeMetric, countGaugePsMetrics)
	metrics := models.Metrics{GaugeMetrics: gaugeMetrics}

	v, err := mem.VirtualMemory()
	if err != nil {
		ag.logger.Errorf("failed to collect virtual mem stats, %s", err)
		errs = append(errs, err)
	} else {
		metrics.GaugeMetrics[totalMemory] = &models.GaugeMetric{
			Name: totalMemory, Type: gauge, Value: float64(v.Total)}
		metrics.GaugeMetrics[freeMemory] = &models.GaugeMetric{
			Name: freeMemory, Type: gauge, Value: float64(v.Free)}
	}

	cpuStat, err := cpu.Percent(0, true)
	if err != nil {
		ag.logger.Errorf("failed to collect virtual cpu stats, %s", err)
		errs = append(errs, err)
	} else {
		metrics.GaugeMetrics[CPUutilization1] = &models.GaugeMetric{
			Name: CPUutilization1, Type: gauge, Value: cpuStat[0]}
	}
	return metrics, errs, maxErrsCount
}
