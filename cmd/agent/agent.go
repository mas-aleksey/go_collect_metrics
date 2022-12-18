package main

import (
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/clients"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"time"
)

var reportInterval = 10 * time.Second
var pollInterval = 2 * time.Second
var baseURL = "http://127.0.0.1:8080"
var timeout = 5 * time.Second

func reportStatistic(statistic *utils.Statistic) {
	metricConfig := clients.NewClientConfig(baseURL, timeout)
	metricClient := clients.NewMetricClient(metricConfig)
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for range ticker.C {
		statCopy := statistic.Copy()
		report := utils.NewJSONReport(statCopy)
		err := metricClient.SendJSONReport(report)
		if err != nil {
			fmt.Println("Fail send report", statCopy.Counter, err)
		} else {
			fmt.Println("Send report successfully", statCopy.Counter)
			statistic.ResetCounter()
		}
	}
}

func updateStatistic(statistic *utils.Statistic) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

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
