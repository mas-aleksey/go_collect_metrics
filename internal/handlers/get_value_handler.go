package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SetValueMetricHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		metric := utils.NewMetric(mType, mName, "0")

		if !metric.IsValidType() {
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
			return
		}
		ok := storage.SetMetricValue(&metric)
		if !ok {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "plain/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(metric.Value))
	}
}
