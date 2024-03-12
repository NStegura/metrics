package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricSort(t *testing.T) {
	gaugeMetrics := make([]GaugeMetric, 0, 3)
	gaugeMetrics = append(gaugeMetrics, GaugeMetric{Name: "C"})
	gaugeMetrics = append(gaugeMetrics, GaugeMetric{Name: "B"})
	gaugeMetrics = append(gaugeMetrics, GaugeMetric{Name: "A"})

	sort.Sort(ByName(gaugeMetrics))

	assert.Equal(t, gaugeMetrics[0].Name, "A")
	assert.Equal(t, gaugeMetrics[1].Name, "B")
	assert.Equal(t, gaugeMetrics[2].Name, "C")
}
