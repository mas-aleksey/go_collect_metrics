package utils

import (
	"strconv"
)

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

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

func IsValidGaugeMetricName(n string) bool {
	switch {
	case contains(n, RuntimeMetricNames):
		return true
	case n == "RandomValue":
		return true
	default:
		return false
	}
}

func IsValidCounterMetricName(n string) bool {
	switch n {
	case "PollCount":
		return true
	default:
		return false
	}
}

func (m Metric) IsValidName() bool {
	switch m.Type {
	case GaugeMetricType:
		return IsValidGaugeMetricName(m.Name)
	case CounterMetricType:
		return IsValidCounterMetricName(m.Name)
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
