package handlers

import (
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
)

func SaveJSONMetricHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metric, err := utils.NewJSONMetric(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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
		w.WriteHeader(http.StatusOK)

	}
}
