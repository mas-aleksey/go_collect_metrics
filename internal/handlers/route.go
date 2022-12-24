package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"time"
)

func GetRouter(memStorage *storage.MemStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(1, "application/json", "text/html", "text/plain"))
	r.Get("/", IndexHandler(memStorage))
	r.Get("/value/{mType}/{mName}", SetValueMetricHandler(memStorage))
	r.Post("/value/", SetValueJSONMetricHandler(memStorage))
	r.Post("/update/{mType}/{mName}/{mValue}", SaveMetricHandler(memStorage))
	r.Post("/update/", SaveJSONMetricHandler(memStorage))
	return r
}
