package handlers

import (
	"encoding/json"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SetValueJSONMetricHandler(db storage.Storage, hashKey string) http.HandlerFunc {
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
		if !metric.IsValidType() {
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
			return
		}
		ok := db.GetBuffer().UpdateJSONMetricValue(&metric)
		metric.Hash = utils.CalcHash(metric.String(), hashKey)
		if !ok {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rest, _ := json.Marshal(metric)
		w.Write(rest)
	}
}
