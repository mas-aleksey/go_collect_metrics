package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"net/http"
)

func GetValueMetricHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		metric, err := db.GetJSONMetric(ctx, mName, mType)
		if err != nil {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(metric.ValueString()))
	}
}
