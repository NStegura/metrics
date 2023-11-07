package models

type MetricName string
type MetricType string

type GaugeMetric struct {
	Name  MetricName
	Type  MetricType
	Value float64
}

type CounterMetric struct {
	Name  MetricName
	Type  MetricType
	Value int64
}

type Metrics struct {
	GaugeMetrics   map[MetricName]GaugeMetric
	CounterMetrics map[MetricName]CounterMetric
}
