package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// SaveBatchJSONMetricHandler - метод для загрузки списка метрик в формате JSON.
// POST /updates/
func SaveBatchJSONMetricHandler(db storage.Storage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := ReadBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("updates request: %s\n", string(body))

		metrics, err := utils.LoadButchJSONMetric(body)
		if err != nil {
			log.Printf("error LoadButchJSONMetric: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, metric := range metrics {
			if err := metric.ValidatesAll(hashKey); err != nil {
				log.Printf("error Validate metric %s: %s", metric.ID, err)
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
		}
		metrics, err = db.UpdateJSONMetrics(ctx, metrics)
		if err != nil {
			log.Printf("error UpdateJSONMetrics: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i, metric := range metrics {
			metrics[i].Hash = utils.CalcHash(metric.String(), hashKey)
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(metrics)
		log.Printf("updates response: %s\n", string(resp))
		w.Write(resp)
	}
}
