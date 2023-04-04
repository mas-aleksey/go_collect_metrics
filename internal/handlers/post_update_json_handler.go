package handlers

import (
	"encoding/json"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

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
		if err := metric.ValidatesAll(hashKey); err != nil {
			switch err {
			case utils.ErrMetricHash:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case utils.ErrMetricType:
				http.Error(w, err.Error(), http.StatusNotImplemented)
			case utils.ErrMetricValue:
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		metric, err = db.UpdateJSONMetric(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metric.Hash = utils.CalcHash(metric.String(), hashKey)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		rest, _ := json.Marshal(metric)
		w.Write(rest)
	}
}
