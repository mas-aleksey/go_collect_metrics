package handlers

import (
	"encoding/json"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SaveBatchJSONMetricHandler(db storage.Storage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ReadBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metrics, err := utils.LoadButchJSONMetric(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, metric := range metrics {
			err = processMetric(&metric, hashKey, db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			metrics[i] = metric
		}
		db.SaveIfSyncMode()
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rest, _ := json.Marshal(metrics)
		w.Write(rest)
	}
}
