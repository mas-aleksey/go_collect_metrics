// Агент для периодической отправки CPU, Memory и других метрик.
package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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
	rateLimit      *int
	pprofMode      *bool
	pprofDuration  *time.Duration
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	reportInterval = flag.Duration("r", 10*time.Second, "report interval")
	pollInterval = flag.Duration("p", 2*time.Second, "pool interval")
	hashKey = flag.String("k", "", "hash key")
	rateLimit = flag.Int("l", 10, "rate limit")
	pprofMode = flag.Bool("pp", false, "pprof mode")
	pprofDuration = flag.Duration("pd", 30*time.Second, "pprof duration")
}

func reportStatistic(statistic *utils.Statistic, config utils.AgentConfig) {
	metricClient := clients.NewMetricClient(config.Address, timeout, config.RateLimit)
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
		statistic.CollectRuntime()
	}
}

func updateMemCPUStatistic(statistic *utils.Statistic, config utils.AgentConfig) {
	ticker := time.NewTicker(config.PollInterval)
	defer ticker.Stop()

	for range ticker.C {
		statistic.CollectMemCPU()
	}
}

func main() {
	flag.Parse()
	config, err := utils.MakeAgentConfig(*address, *reportInterval, *pollInterval, *hashKey, *rateLimit)
	if err != nil {
		log.Fatal(err)
	}
	stat := utils.NewStatistic()
	go updateStatistic(stat, config)
	go updateMemCPUStatistic(stat, config)
	go reportStatistic(stat, config)

	if *pprofMode {
		time.Sleep(*pprofDuration)
		fmem, err := os.Create(`mem.pprof`)
		if err != nil {
			panic(err)
		}
		defer fmem.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fmem); err != nil {
			panic(err)
		}
	} else {
		select {}
	}
}
