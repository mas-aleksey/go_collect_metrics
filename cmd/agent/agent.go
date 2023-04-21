package main

import (
	"flag"
	"github.com/tiraill/go_collect_metrics/internal/clients"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"time"
)

var (
	address        *string
	reportInterval *time.Duration
	pollInterval   *time.Duration
	timeout        = 5 * time.Second
	hashKey        *string
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	reportInterval = flag.Duration("r", 10*time.Second, "report interval")
	pollInterval = flag.Duration("p", 2*time.Second, "pool interval")
	hashKey = flag.String("k", "", "hash key")
}

func reportStatistic(statistic *utils.Statistic, config utils.AgentConfig) {
	metricClient := clients.NewMetricClient(config.Address, timeout)
	ticker := time.NewTicker(config.ReportInterval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Sending report...")
		statCopy := statistic.Copy()
		report := utils.NewJSONReport(statCopy, config.HashKey)
		err := metricClient.SendBatchJSONReport(report)
		if err != nil {
			log.Println("Fail send report", statCopy.Counter, err)
		} else {
			log.Println("Send report successfully", statCopy.Counter)
			statistic.ResetCounter()
		}
	}
}

func updateStatistic(statistic *utils.Statistic, config utils.AgentConfig) {
	ticker := time.NewTicker(config.PollInterval)
	defer ticker.Stop()

	for range ticker.C {
		statistic.Collect()
	}
}

func main() {
	flag.Parse()
	config := utils.MakeAgentConfig(*address, *reportInterval, *pollInterval, *hashKey)
	stat := utils.NewStatistic()
	go updateStatistic(stat, config)
	go reportStatistic(stat, config)
	select {}
}
