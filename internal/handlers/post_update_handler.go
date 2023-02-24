package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
)

func SaveMetricHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")
		metric, err := utils.NewJSONMetric(mType, mName, mValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		db.GetBuffer().PutJSONMetric(metric)
		db.SaveIfSyncMode()
		w.WriteHeader(http.StatusOK)

	}
}
