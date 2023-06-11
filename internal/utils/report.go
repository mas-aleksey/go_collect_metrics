package utils

import (
	"reflect"
)

var ReportCount = 32

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

	totalMemoryMetric := NewGaugeJSONMetric("TotalMemory", ToFloat64(statistic.MemStat.Total))
	totalMemoryMetric.Hash = CalcHash(totalMemoryMetric.String(), hashKey)
	metrics = append(metrics, totalMemoryMetric)

	freeMemoryMetric := NewGaugeJSONMetric("FreeMemory", ToFloat64(statistic.MemStat.Free))
	freeMemoryMetric.Hash = CalcHash(freeMemoryMetric.String(), hashKey)
	metrics = append(metrics, freeMemoryMetric)

	cpuUtilization1 := NewGaugeJSONMetric("CPUutilization1", ToFloat64(statistic.CPUutilization1[0]))
	cpuUtilization1.Hash = CalcHash(cpuUtilization1.String(), hashKey)
	metrics = append(metrics, cpuUtilization1)

	for _, metricName := range RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		metric := NewGaugeJSONMetric(metricName, ToFloat64(val))
		metric.Hash = CalcHash(metric.String(), hashKey)
		metrics = append(metrics, metric)
	}
	return &JSONReport{metrics}
}
