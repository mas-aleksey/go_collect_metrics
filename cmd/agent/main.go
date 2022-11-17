package main

import (
	"clients"
	"time"
	"utils"
)

var reportInterval = 10 * time.Second
var pollInterval = 2 * time.Second

func reportStatistic(statistic *utils.Statistic) {
	metricClient := clients.NewMetricClient("http://127.0.0.1:8080")
	ticker := time.NewTicker(reportInterval)
	for _ = range ticker.C {
		metricClient.SendMetrics(*statistic)
	}
}

func updateStatistic(statistic *utils.Statistic) {
	ticker := time.NewTicker(pollInterval)
	for _ = range ticker.C {
		statistic.Collect()
	}
}

func main() {
	stat := utils.NewStatistic()
	go updateStatistic(stat)
	go reportStatistic(stat)
	select {}
}
