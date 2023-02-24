package handlers

import (
	"encoding/json"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"net/http"
)

func SaveBatchJSONMetricHandler(db storage.Storage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//r.Context()
		body, err := ReadBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("updates request: %s\n", string(body))
		log.Printf("hashKey: %s\n", hashKey)

		metrics, err := utils.LoadButchJSONMetric(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, metric := range metrics {
			if err := metric.ValidatesAll(hashKey); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		metrics, err = db.UpdateJSONMetrics(metrics)
		if err != nil {
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
