package models

import (
	"strconv"
)

type Metric struct {
	Name  string
	Type  string
	Value string
}

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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func CastToGauge(m Metric) (GaugeMetric, error) {
	value, err := strconv.ParseFloat(m.Value, 64)
	if err != nil {
		return GaugeMetric{}, err
	}
	return GaugeMetric{m.Name, m.Type, value}, nil
}

func CastToCounter(m Metric) (CounterMetric, error) {
	value, err := strconv.ParseInt(m.Value, 10, 64)
	if err != nil {
		return CounterMetric{}, err
	}
	return CounterMetric{m.Name, m.Type, value}, nil
}
