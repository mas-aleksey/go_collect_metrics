package utils

import (
	"encoding/json"
	"fmt"
)

type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  *string  `json:"hash,omitempty"`  // значение хеш-функции
}

func NewCounterJSONMetric(mName string, delta int64, hashKey string) JSONMetric {
	m := JSONMetric{
		ID:    mName,
		MType: string(CounterMetricType),
		Delta: &delta,
	}
	m.Hash = CalcHash(m.String(), hashKey)
	return m
}

func NewGaugeJSONMetric(mName string, value float64, hashKey string) JSONMetric {
	m := JSONMetric{
		ID:    mName,
		MType: string(GaugeMetricType),
		Value: &value,
	}
	m.Hash = CalcHash(m.String(), hashKey)
	return m
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
