package utils

import (
	"reflect"
)

var ReportCount = 29

type Report struct {
	Metrics []JSONMetric
}

type JSONReport struct {
	Metrics []JSONMetric
}

func NewJSONReport(statistic *Statistic, hashKey string) *JSONReport {
	metrics := make([]JSONMetric, 0, ReportCount)
	metrics = append(metrics, NewCounterJSONMetric("PollCount", statistic.Counter, hashKey))
	metrics = append(metrics, NewGaugeJSONMetric("RandomValue", statistic.RndValue, hashKey))
	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metrics = append(metrics, NewGaugeJSONMetric(metricName, ToFloat64(val), hashKey))
	}
	return &JSONReport{metrics}
}
