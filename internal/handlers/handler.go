package handlers

import (
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
	"strings"
)

func SaveMetricHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request:", r.Method, r.URL)
		switch r.Method {
		case "POST":
			fragments := strings.Split(r.URL.Path, "/")
			if len(fragments) != 5 {
				errMsg := "Wrong path. Expected: 'HOST/update/metric_type/metric_name/metric_value'"
				http.Error(w, errMsg, http.StatusNotFound)
				return
			}
			metric := utils.NewMetric(fragments[2], fragments[3], fragments[4])

			if !metric.IsValidType() {
				http.Error(w, "Invalid metric type", http.StatusNotImplemented)
				return
			}
			if !metric.IsValidValue() {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			storage.SaveMetric(metric)
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method "+r.Method+" not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
