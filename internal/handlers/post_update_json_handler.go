package handlers

import (
	"encoding/json"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SaveJSONMetricHandler(storage *storage.MemStorage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ReadBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metric, err := utils.LoadJSONMetric(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if !metric.IsValidHash(hashKey) {
			http.Error(w, "Invalid metric hash", http.StatusBadRequest)
			return
		}
		if !metric.IsValidType() {
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
			return
		}
		if !metric.IsValidValue() {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		storage.SaveJSONMetric(metric)
		storage.SaveToFileIfSyncMode()
		storage.SetJSONMetricValue(&metric)
		metric.Hash = utils.CalcHash(metric.String(), hashKey)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rest, _ := json.Marshal(metric)
		w.Write(rest)
	}
}
