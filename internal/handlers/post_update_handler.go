package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SaveMetricHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")
		metric := utils.NewMetric(mType, mName, mValue)

		if !metric.IsValidType() {
			http.Error(w, "Invalid metric type", http.StatusNotImplemented)
			return
		}
		if !metric.IsValidValue() {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		storage.SaveMetric(metric)
		storage.SaveToFileIfSyncMode()
		w.WriteHeader(http.StatusOK)

	}
}
