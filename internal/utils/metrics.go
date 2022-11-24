package utils

import (
	"strconv"
)

type MetricType string

const (
	GaugeMetricType   MetricType = "gauge"
	CounterMetricType MetricType = "counter"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value string
}

func NewMetric(metricType string, metricName string, metricValue string) Metric {
	return Metric{
		Type:  MetricType(metricType),
		Name:  metricName,
		Value: metricValue,
	}
}

func (m Metric) IsValidType() bool {
	switch m.Type {
	case GaugeMetricType, CounterMetricType:
		return true
	default:
		return false
	}
}

func (m Metric) IsValidValue() bool {
	switch m.Type {
	case GaugeMetricType: // float64
		if _, err := strconv.ParseFloat(m.Value, 64); err == nil {
			return true
		}
		return false
	case CounterMetricType: // int64
		if _, err := strconv.ParseInt(m.Value, 10, 64); err == nil {
			return true
		}
		return false
	default:
		return false
	}
}
