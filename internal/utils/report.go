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

	pollCountMetric := NewCounterJSONMetric("PollCount", statistic.Counter)
	pollCountMetric.Hash = CalcHash(pollCountMetric.String(), hashKey)
	metrics = append(metrics, pollCountMetric)

	randomValueMetric := NewGaugeJSONMetric("RandomValue", statistic.RndValue)
	randomValueMetric.Hash = CalcHash(randomValueMetric.String(), hashKey)
	metrics = append(metrics, randomValueMetric)

	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metric := NewGaugeJSONMetric(metricName, ToFloat64(val))
		metric.Hash = CalcHash(metric.String(), hashKey)
		metrics = append(metrics, NewGaugeJSONMetric(metricName, ToFloat64(val)))
	}
	return &JSONReport{metrics}
}
