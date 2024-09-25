package metric

import "github.com/NStegura/metrics/internal/app/agent/models"

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func CastToMetrics(m models.Metrics) (metrics []Metrics) {
	for _, metric := range m.CounterMetrics {
		metrics = append(metrics, Metrics{
			ID:    string(metric.Name),
			MType: string(metric.Type),
			Delta: &metric.Value,
		})
	}
	for _, metric := range m.GaugeMetrics {
		metrics = append(metrics, Metrics{
			ID:    string(metric.Name),
			MType: string(metric.Type),
			Value: &metric.Value,
		})
	}
	return metrics
}
