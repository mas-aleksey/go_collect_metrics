package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"os"
	"sync"
	"time"
)

type MemStorage struct {
	GaugeMetrics   map[string]float64   `json:"GaugeMetrics"`
	CounterMetrics map[string]int64     `json:"CounterMetrics"`
	Mutex          sync.RWMutex         `json:"-"`
	Config         *utils.StorageConfig `json:"-"`
}

func flushBackground(m *MemStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		m.saveToFile()
	}
}

func (m *MemStorage) Init(ctx context.Context) error {
	if m.Config.StoreInterval != 0 {
		go flushBackground(m, m.Config.StoreInterval)
	}
	if !m.Config.Restore {
		return fmt.Errorf("no need restore")
	}
	if m.Config.StoreFile == "" {
		return fmt.Errorf("filename is empty")
	}
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	data, err := os.ReadFile(m.Config.StoreFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) Close(ctx context.Context) {
	m.saveToFile()
}

func (m *MemStorage) Ping(ctx context.Context) bool {
	return true
}

func (m *MemStorage) UpdateJSONMetric(ctx context.Context, metricIn utils.JSONMetric) (utils.JSONMetric, error) {
	metricOut := m.updateJSONMetric(metricIn)
	if m.Config.StoreInterval == 0 {
		m.saveToFile()
	}
	return metricOut, nil
}

func (m *MemStorage) updateJSONMetric(metricIn utils.JSONMetric) utils.JSONMetric {
	metricOut := utils.JSONMetric{
		ID:    metricIn.ID,
		MType: metricIn.MType,
	}
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	switch metricIn.MType {
	case "gauge":
		val := *metricIn.Value
		m.GaugeMetrics[metricIn.ID] = val
		metricOut.Value = &val
	case "counter":
		val := *metricIn.Delta + m.CounterMetrics[metricIn.ID]
		m.CounterMetrics[metricIn.ID] = val
		metricOut.Delta = &val
	}
	return metricOut
}

func (m *MemStorage) UpdateJSONMetrics(ctx context.Context, metricsIn []utils.JSONMetric) ([]utils.JSONMetric, error) {
	metricsOut := make([]utils.JSONMetric, 0)
	for _, metricIn := range metricsIn {
		metricOut := m.updateJSONMetric(metricIn)
		metricsOut = append(metricsOut, metricOut)
	}
	if m.Config.StoreInterval == 0 {
		m.saveToFile()
	}
	return metricsOut, nil
}

func (m *MemStorage) GetJSONMetric(ctx context.Context, mName, mType string) (utils.JSONMetric, error) {
	metric := utils.JSONMetric{
		ID:    mName,
		MType: mType,
	}
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	switch mType {
	case "gauge":
		val, ok := m.GaugeMetrics[metric.ID]
		if !ok {
			return metric, fmt.Errorf("gauge metric no found")
		}
		metric.Value = &val
	case "counter":
		val, ok := m.CounterMetrics[metric.ID]
		if !ok {
			return metric, fmt.Errorf("counter metric no found")
		}
		metric.Delta = &val
	default:
		return metric, fmt.Errorf("invalid metric type")
	}
	return metric, nil
}

func (m *MemStorage) GetAllMetrics(ctx context.Context) ([]utils.JSONMetric, error) {
	metrics := make([]utils.JSONMetric, 0)
	for name, val := range m.GaugeMetrics {
		metrics = append(metrics, utils.NewGaugeJSONMetric(name, val))
	}
	for name, val := range m.CounterMetrics {
		metrics = append(metrics, utils.NewCounterJSONMetric(name, val))
	}
	return metrics, nil
}

func (m *MemStorage) saveToFile() {
	if m.Config.StoreFile == "" {
		log.Print("Failed save to file: filename is empty")
	}
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	file, err := os.Create(m.Config.StoreFile)

	if err != nil {
		log.Print("Failed save to file", err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Print("Failed save to file", err)
	}
	_, err = file.Write(data)
	if err != nil {
		log.Print("Failed save to file", err)
	}
	err = file.Close()
	if err != nil {
		log.Print("Failed save to file", err)
	}
	log.Print("Save storage to file")
}
