package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// SaveMetricHandler - метод для загрузки метрики.
// POST /update/{mType}/{mName}/{mValue}.
func SaveMetricHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")
		metric, err := utils.NewJSONMetric(mType, mName, mValue)
		if err != nil {
			switch err {
			case utils.ErrMetricType:
				http.Error(w, err.Error(), http.StatusNotImplemented)
			case utils.ErrMetricValue:
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		_, err = db.UpdateJSONMetric(ctx, metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
