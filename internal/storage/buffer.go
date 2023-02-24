package storage

import (
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"sync"
)

type Buffer struct {
	GaugeMetrics   map[string]float64 `json:"GaugeMetrics"`
	CounterMetrics map[string]int64   `json:"CounterMetrics"`
	Mutex          sync.RWMutex       `json:"-"`
}

func NewBuffer() *Buffer {
	return &Buffer{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
}

func (b *Buffer) PutJSONMetric(metrics utils.JSONMetric) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	switch metrics.MType {
	case "gauge":
		b.GaugeMetrics[metrics.ID] = *metrics.Value
	case "counter":
		b.CounterMetrics[metrics.ID] += *metrics.Delta
	}
}

func (b *Buffer) UpdateJSONMetricValue(metric *utils.JSONMetric) bool {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	switch metric.MType {
	case "gauge":
		val, ok := b.GaugeMetrics[metric.ID]
		metric.Value = &val
		return ok
	case "counter":
		val, ok := b.CounterMetrics[metric.ID]
		metric.Delta = &val
		return ok
	default:
		return false
	}
}
