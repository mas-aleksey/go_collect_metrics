package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
	"utils"
)

var baseUrl = "http://127.0.0.1:8080"
var reportInterval = 10 * time.Second
var pollInterval = 2 * time.Second

func postMetric(client *http.Client, mType string, mName string, mValue string) {
	endpoint, _ := url.JoinPath(baseUrl, "update", mType, mName, mValue)
	response, err := client.Post(endpoint, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Response:", response.Status, strings.TrimSpace(string(body)))
}

func sendMetrics(client *http.Client, statistic utils.Statistic) {
	postMetric(client, "counter", "PollCount", utils.ToStr(statistic.Counter))
	for _, metricName := range utils.RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		postMetric(client, "gauge", metricName, utils.ToStr(val))
	}
	postMetric(client, "gauge", "RandomValue", utils.ToStr(statistic.RndValue))
}

func reportStatistic(statistic *utils.Statistic) {
	client := &http.Client{}
	ticker := time.NewTicker(reportInterval)
	for _ = range ticker.C {
		sendMetrics(client, *statistic)
		fmt.Println("Report statistic", statistic.Counter)
	}
}

func updateStatistic(statistic *utils.Statistic) {
	ticker := time.NewTicker(pollInterval)
	for _ = range ticker.C {
		statistic.Collect()
		fmt.Println("Collect statistic", statistic.Counter)
	}
}

func main() {
	stat := utils.NewStatistic()
	go updateStatistic(stat)
	go reportStatistic(stat)
	select {}
}
