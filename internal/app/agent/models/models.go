package models

type MetricName string
type MetricType string

type GaugeMetric struct {
	Name  MetricName `json:"id"`
	Type  MetricType `json:"type"`
	Value float64    `json:"value"`
}

type CounterMetric struct {
	Name  MetricName `json:"id"`
	Type  MetricType `json:"type"`
	Value int64      `json:"delta"`
}

// Metrics - отправляемые метрики.
type Metrics struct {
	GaugeMetrics   map[MetricName]*GaugeMetric
	CounterMetrics map[MetricName]*CounterMetric
}
