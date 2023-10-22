package utils

import (
	"fmt"
	"reflect"
)

var ReportCount = 47

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

	totalMemoryMetric := NewGaugeJSONMetric("TotalMemory", ToFloat64(statistic.MemStat.Total))
	totalMemoryMetric.Hash = CalcHash(totalMemoryMetric.String(), hashKey)
	metrics = append(metrics, totalMemoryMetric)

	freeMemoryMetric := NewGaugeJSONMetric("FreeMemory", ToFloat64(statistic.MemStat.Free))
	freeMemoryMetric.Hash = CalcHash(freeMemoryMetric.String(), hashKey)
	metrics = append(metrics, freeMemoryMetric)

	for i, utilization := range statistic.CPUUtilization {
		metricName := fmt.Sprintf("CPUutilization%d", i+1)
		cpuUtilization1 := NewGaugeJSONMetric(metricName, ToFloat64(utilization))
		cpuUtilization1.Hash = CalcHash(cpuUtilization1.String(), hashKey)
		metrics = append(metrics, cpuUtilization1)
	}
	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metric := NewGaugeJSONMetric(metricName, ToFloat64(val))
		metric.Hash = CalcHash(metric.String(), hashKey)
		metrics = append(metrics, metric)
	}
	return &JSONReport{metrics}
}
