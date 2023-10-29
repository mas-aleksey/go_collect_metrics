package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

// ErrMetricHash ошибка невалидного хеша метрики.
var ErrMetricHash = errors.New("invalid metric hash")

// ErrMetricType ошибка невалидного типа метрики.
var ErrMetricType = errors.New("invalid metric type")

// ErrMetricValue ошибка невалидного значения метрики.
var ErrMetricValue = errors.New("invalid metric value")

// JSONMetric - структура метрики в формате JSON.
type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  *string  `json:"hash,omitempty"`  // значение хеш-функции
}

// NewCounterJSONMetric - метод создания объекта метрики с типом counter.
func NewCounterJSONMetric(mName string, delta int64) JSONMetric {
	return JSONMetric{
		ID:    mName,
		MType: "counter",
		Delta: &delta,
	}
}

// NewGaugeJSONMetric - метод создания объекта метрики с типом gauge.
func NewGaugeJSONMetric(mName string, value float64) JSONMetric {
	return JSONMetric{
		ID:    mName,
		MType: "gauge",
		Value: &value,
	}
}

// NewJSONMetric - общий метод создания объекта метрики.
func NewJSONMetric(metricType, metricName, metricValue string) (JSONMetric, error) {
	m := JSONMetric{}
	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return m, ErrMetricValue
		}
		return NewGaugeJSONMetric(metricName, val), nil
	case "counter":
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return m, ErrMetricValue
		}
		return NewCounterJSONMetric(metricName, val), nil
	default:
		return m, ErrMetricType
	}
}

// LoadJSONMetric - метод парсинга метрики из данных в формате JSON.
func LoadJSONMetric(body []byte) (JSONMetric, error) {
	var metric JSONMetric
	if err := json.Unmarshal(body, &metric); err != nil {
		return metric, err
	}
	return metric, nil
}

// LoadButchJSONMetric - метод парсинга списка метрик из данных в формате JSON.
func LoadButchJSONMetric(body []byte) ([]JSONMetric, error) {
	var metrics []JSONMetric
	if err := json.Unmarshal(body, &metrics); err != nil {
		return metrics, err
	}
	return metrics, nil
}

// Метод приведения JSONMetric к строке.
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

// ValueString - метод приведения значения метрики к строке.
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

// IsValidType - метод валидации типа метрики.
func (m JSONMetric) IsValidType() bool {
	switch m.MType {
	case "gauge", "counter":
		return true
	default:
		return false
	}
}

// IsValidValue - метод валидации значения метрики.
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

// IsValidHash - метод валидации хеш-суммы метрики.
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

// ValidatesAll - общий метод валидации метрики.
func (m JSONMetric) ValidatesAll(hashKey string) error {
	if !m.IsValidHash(hashKey) {
		return ErrMetricHash
	}
	if !m.IsValidType() {
		return ErrMetricType
	}
	if !m.IsValidValue() {
		return ErrMetricValue
	}
	return nil
}
