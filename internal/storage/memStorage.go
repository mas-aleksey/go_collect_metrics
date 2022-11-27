package storage

import (
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"strconv"
)

type MemStorage struct {
	GaugeMetrics   map[string]float64
	CounterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (m *MemStorage) SaveMetric(metric utils.Metric) {
	switch metric.Type {
	case utils.GaugeMetricType:
		val, _ := strconv.ParseFloat(metric.Value, 64)
		m.GaugeMetrics[metric.Name] = val
	case utils.CounterMetricType:
		val, _ := strconv.ParseInt(metric.Value, 10, 64)
		m.CounterMetrics[metric.Name] += val
	}
}

func (m *MemStorage) SaveJSONMetric(metrics utils.JSONMetric) {
	switch metrics.MType {
	case "gauge":
		m.GaugeMetrics[metrics.ID] = *metrics.Value
	case "counter":
		m.CounterMetrics[metrics.ID] += *metrics.Delta
	}
}

func (m *MemStorage) SetMetricValue(metric *utils.Metric) bool {
	switch metric.Type {
	case utils.GaugeMetricType:
		val, ok := m.GaugeMetrics[metric.Name]
		metric.Value = utils.ToStr(val)
		return ok
	case utils.CounterMetricType:
		val, ok := m.CounterMetrics[metric.Name]
		metric.Value = utils.ToStr(val)
		return ok
	default:
		return false
	}
}

func (m *MemStorage) SetJSONMetricValue(metric *utils.JSONMetric) bool {
	switch metric.MType {
	case "gauge":
		val, ok := m.GaugeMetrics[metric.ID]
		metric.Value = &val
		return ok
	case "counter":
		val, ok := m.CounterMetrics[metric.ID]
		metric.Delta = &val
		return ok
	default:
		return false
	}
}
