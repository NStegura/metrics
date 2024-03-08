package agent

import "github.com/NStegura/metrics/internal/clients/metric"

// MetricCli определяет интерфейс для работы с метриками.
type MetricCli interface {
	UpdateMetrics([]metric.Metrics) error
}
