package storage

import (
	"strconv"
	"utils"
)

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string][]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string][]int64),
	}
}

func (m *MemStorage) SaveMetric(metric utils.Metric) {
	switch metric.Type {
	case utils.GaugeMetricType:
		val, _ := strconv.ParseFloat(metric.Value, 64)
		m.GaugeMetrics[metric.Name] = val
	case utils.CounterMetricType:
		val, _ := strconv.ParseInt(metric.Value, 10, 64)
		m.CounterMetrics[metric.Name] = append(m.CounterMetrics[metric.Name], val)
	}
}
