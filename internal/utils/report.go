package utils

import (
	"reflect"
)

var ReportCount = 29

type Report struct {
	Metrics []Metric
}

type JSONReport struct {
	Metrics []JSONMetric
}

func NewJSONReport(statistic *Statistic) *JSONReport {
	metrics := make([]JSONMetric, 0, ReportCount)
	metrics = append(metrics, NewCounterJSONMetric("PollCount", statistic.Counter))
	metrics = append(metrics, NewGaugeJSONMetric("RandomValue", statistic.RndValue))
	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metrics = append(metrics, NewGaugeJSONMetric(metricName, ToFloat64(val)))
	}
	return &JSONReport{metrics}
}
