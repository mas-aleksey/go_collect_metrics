package main

import (
	"github.com/tiraill/go_collect_metrics/internal/clients"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"time"
)

var reportInterval = 10 * time.Second
var pollInterval = 2 * time.Second

func reportStatistic(statistic *utils.Statistic) {
	metricClient := clients.NewMetricClient("http://127.0.0.1:8080")
	ticker := time.NewTicker(reportInterval)
	for range ticker.C {
		metricClient.SendMetrics(*statistic)
	}
}

func updateStatistic(statistic *utils.Statistic) {
	ticker := time.NewTicker(pollInterval)
	for range ticker.C {
		statistic.Collect()
	}
}

func main() {
	stat := utils.NewStatistic()
	go updateStatistic(stat)
	go reportStatistic(stat)
	select {}
}
