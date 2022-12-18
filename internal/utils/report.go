package utils

import (
	"reflect"
)

var ReportCount = 29

type Report struct {
	Metrics []Metric
}

func NewReport(statistic *Statistic) *Report {
	metrics := make([]Metric, 0, ReportCount)
	metrics = append(metrics, NewMetric("counter", "PollCount", ToStr(statistic.Counter)))
	metrics = append(metrics, NewMetric("gauge", "RandomValue", ToStr(statistic.RndValue)))
	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metrics = append(metrics, NewMetric("gauge", metricName, ToStr(val)))
	}
	return &Report{metrics}
}
