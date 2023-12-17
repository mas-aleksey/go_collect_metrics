package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// GetJSONMetricHandler - метод получения значения метрики в формате JSON
// POST /value/
func GetJSONMetricHandler(db storage.Storage, hashKey string, privateKey *utils.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := ReadEncryptedBody(r, privateKey)
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
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
		metric, err = db.GetJSONMetric(ctx, metric.ID, metric.MType)
		if err != nil {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		metric.Hash = utils.CalcHash(metric.String(), hashKey)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(metric)
		log.Printf("post values response: %s\n", string(resp))
		w.Write(resp)
	}
}
