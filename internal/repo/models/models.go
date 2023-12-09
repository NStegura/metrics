package models

type GaugeMetric struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type CounterMetric struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value int64  `json:"value"`
}
