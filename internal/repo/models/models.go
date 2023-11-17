package models

type GaugeMetric struct {
	Name  string
	Type  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Type  string
	Value int64
}
