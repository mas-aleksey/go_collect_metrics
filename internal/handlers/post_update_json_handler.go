package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func processMetric(metric *utils.JSONMetric, hashKey string, db storage.Storage) error {
	if !metric.IsValidHash(hashKey) {
		return fmt.Errorf("invalid metric hash")
	}
	if !metric.IsValidType() {
		return fmt.Errorf("invalid metric type")
	}
	if !metric.IsValidValue() {
		return fmt.Errorf("invalid metric value")
	}
	db.GetBuffer().PutJSONMetric(*metric)
	db.GetBuffer().UpdateJSONMetricValue(metric)
	metric.Hash = utils.CalcHash(metric.String(), hashKey)
	return nil
}

func SaveJSONMetricHandler(db storage.Storage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ReadBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metric, err := utils.LoadJSONMetric(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = processMetric(&metric, hashKey, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		db.SaveIfSyncMode()
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rest, _ := json.Marshal(metric)
		w.Write(rest)
	}
}
