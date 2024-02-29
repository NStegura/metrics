package agent

import "github.com/NStegura/metrics/internal/clients/metric"

type MetricCli interface {
	UpdateMetrics([]metric.Metrics) error
}
