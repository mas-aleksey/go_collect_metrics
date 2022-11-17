package storage

import (
	"strconv"
	"utils"
)

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics []int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make([]int64, 0),
	}
}

func (m *MemStorage) SaveMetric(metric utils.Metric) {
	switch metric.Type {
	case utils.GaugeMetricType:
		val, _ := strconv.ParseFloat(metric.Value, 64)
		m.GaugeMetrics[metric.Name] = val
	case utils.CounterMetricType:
		val, _ := strconv.ParseInt(metric.Value, 10, 64)
		m.CounterMetrics = append(m.CounterMetrics, val)
	}
}
