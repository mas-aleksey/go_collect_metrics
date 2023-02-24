package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  *string  `json:"hash,omitempty"`  // значение хеш-функции
}

func NewCounterJSONMetric(mName string, delta int64) JSONMetric {
	return JSONMetric{
		ID:    mName,
		MType: "counter",
		Delta: &delta,
	}
}

func NewGaugeJSONMetric(mName string, value float64) JSONMetric {
	return JSONMetric{
		ID:    mName,
		MType: "gauge",
		Value: &value,
	}
}

func NewJSONMetric(metricType, metricName, metricValue string) (JSONMetric, error) {
	m := JSONMetric{}
	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return m, fmt.Errorf("invalid gauge metric value")
		}
		return NewGaugeJSONMetric(metricName, val), nil
	case "counter":
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return m, fmt.Errorf("invalid counter metric value")
		}
		return NewCounterJSONMetric(metricName, val), nil
	default:
		return m, fmt.Errorf("invalid metric type")
	}
}

func LoadJSONMetric(body []byte) (JSONMetric, error) {
	var metric JSONMetric
	if err := json.Unmarshal(body, &metric); err != nil {
		return metric, err
	}
	return metric, nil
}

func LoadButchJSONMetric(body []byte) ([]JSONMetric, error) {
	var metrics []JSONMetric
	if err := json.Unmarshal(body, &metrics); err != nil {
		return metrics, err
	}
	return metrics, nil
}

func (m JSONMetric) String() string {
	switch m.MType {
	case "gauge":
		return fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		return fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	default:
		return ""
	}
}

func (m JSONMetric) ValueString() string {
	switch m.MType {
	case "gauge":
		return fmt.Sprintf("%g", *m.Value)
	case "counter":
		return fmt.Sprintf("%d", *m.Delta)
	default:
		return ""
	}
}

func (m JSONMetric) IsValidType() bool {
	switch m.MType {
	case "gauge", "counter":
		return true
	default:
		return false
	}
}

func (m JSONMetric) IsValidValue() bool {
	switch m.MType {
	case "gauge":
		return m.Value != nil
	case "counter":
		return m.Delta != nil
	default:
		return false
	}
}

func (m JSONMetric) IsValidHash(hashKey string) bool {
	if m.Hash == nil {
		return true
	}
	if hashKey == "" {
		return true
	}
	actualHash := CalcHash(m.String(), hashKey)
	return *actualHash == *m.Hash
}

func (m JSONMetric) ValidatesAll(hashKey string) error {
	if !m.IsValidHash(hashKey) {
		return fmt.Errorf("invalid metric hash")
	}
	if !m.IsValidType() {
		return fmt.Errorf("invalid metric type")
	}
	if !m.IsValidValue() {
		return fmt.Errorf("invalid metric value")
	}
	return nil
}
