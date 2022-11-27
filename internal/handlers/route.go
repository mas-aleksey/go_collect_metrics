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
	r.Get("/", IndexHandler(memStorage))
	r.Get("/value/{mType}/{mName}", GetValueMetricHandler(memStorage))
	r.Post("/update/{mType}/{mName}/{mValue}", SaveMetricHandler(memStorage))
	r.Post("/update", SaveJsonMetricHandler(memStorage))
	return r
}
