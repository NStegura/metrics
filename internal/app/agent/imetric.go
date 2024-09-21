package agent

import (
	"context"

	"github.com/NStegura/metrics/internal/clients/metric"
)

// MetricCli определяет интерфейс для работы с метриками.
type MetricCli interface {
	UpdateMetrics(context.Context, []metric.Metrics) error
}
