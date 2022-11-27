package utils

import (
	"encoding/json"
)

type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewJSONMetric(body []byte) (JSONMetric, error) {
	metric := JSONMetric{}
	if err := json.Unmarshal(body, &metric); err != nil {
		return metric, err
	}
	return metric, nil
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
