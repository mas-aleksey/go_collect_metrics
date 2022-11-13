package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"utils"
)

func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:", r.Method, r.URL)
	switch r.Method {
	case "POST":
		fragments := strings.Split(r.URL.Path, "/")
		if len(fragments) != 5 {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Wrong path. Expected: 'HOST/update/metric_type/metric_name/metric_value'")
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Expected 'Content-Type' only 'text/plain'")
			return
		}
		metric := utils.NewMetric(fragments[2], fragments[3], fragments[4])

		if !metric.IsValidType() {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Wrong metric type. Expected: [", utils.GaugeMetricType, utils.CounterMetricType, "]")
			return
		}
		if !metric.IsValidName() {
			w.WriteHeader(400)
			fmt.Fprintln(w, "For metric type:", metric.Type, "got wrong metric name:", metric.Name)
			return
		}
		if !metric.IsValidValue() {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Invalid metric value")
			return
		}
		w.WriteHeader(201)
		fmt.Fprintln(w, "Save metric type:", metric.Type, "name:", metric.Name, "value:", metric.Value)
	default:
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method", r.Method, "not allowed")
	}
}

func main() {
	http.HandleFunc("/update/", updateMetricHandler)
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
