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

type ByName []GaugeMetric

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
