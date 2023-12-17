// Агент для периодической отправки CPU, Memory и других метрик.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tiraill/go_collect_metrics/internal/clients"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

var (
	address        *string
	reportInterval *time.Duration
	pollInterval   *time.Duration
	timeout        = 5 * time.Second
	hashKey        *string
	cryptoKey      *string
	configFile     *string
	rateLimit      *int
	buildVersion   = "N/A"
	buildDate      = "N/A"
	buildCommit    = "N/A"
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	reportInterval = flag.Duration("r", 10*time.Second, "report interval")
	pollInterval = flag.Duration("p", 2*time.Second, "pool interval")
	hashKey = flag.String("k", "", "hash key")
	cryptoKey = flag.String("crypto-key", "", "public crypto key")
	configFile = flag.String("config", "", "config file")
	rateLimit = flag.Int("l", 10, "rate limit")
}

func reportStatistic(statistic *utils.Statistic, config utils.AgentConfig, metricClient *clients.MetricClient) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
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

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	flag.Parse()
	config, err := utils.MakeAgentConfig(*configFile, *address, *reportInterval, *pollInterval, *hashKey, *cryptoKey, *rateLimit)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	stat := utils.NewStatistic()
	metricClient := clients.NewMetricClient(config.Address, timeout, config.RateLimit, config.CryptoKey)

	reportStatisticTicker := time.NewTicker(config.ReportInterval)
	updateStatisticTicker := time.NewTicker(config.PollInterval)
	updateMemCPUStatisticTicker := time.NewTicker(config.PollInterval)

	log.Print("Agent Started")
	for {
		select {
		case <-reportStatisticTicker.C:
			reportStatistic(stat, config, metricClient)
		case <-updateStatisticTicker.C:
			stat.CollectRuntime()
		case <-updateMemCPUStatisticTicker.C:
			stat.CollectMemCPU()
		case s := <-done:
			log.Print("Agent Stopped. Signal: ", s)
			reportStatisticTicker.Stop()
			updateStatisticTicker.Stop()
			updateMemCPUStatisticTicker.Stop()
			reportStatistic(stat, config, metricClient)
			log.Print("Exit")
			return
		}
	}
}
