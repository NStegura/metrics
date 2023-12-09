package models

import (
	"fmt"
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
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func CastToGauge(m Metric) (GaugeMetric, error) {
	value, err := strconv.ParseFloat(m.Value, 64)
	if err != nil {
		return GaugeMetric{}, fmt.Errorf("parse float failed: %w", err)
	}
	return GaugeMetric{m.Name, m.Type, value}, nil
}

func CastToCounter(m Metric) (CounterMetric, error) {
	value, err := strconv.ParseInt(m.Value, 10, 64)
	if err != nil {
		return CounterMetric{}, fmt.Errorf("parse int failed: %w", err)
	}
	return CounterMetric{m.Name, m.Type, value}, nil
}
