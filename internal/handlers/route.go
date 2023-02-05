package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"time"
)

func GetRouter(db storage.Storage, config utils.ServerConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(1, "application/json", "text/html", "text/plain"))
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Get("/", IndexHandler(db))
	r.Get("/ping", GetPingHandler(db))
	r.Get("/value/{mType}/{mName}", SetValueMetricHandler(db))
	r.Post("/value/", SetValueJSONMetricHandler(db, config.HashKey))
	r.Post("/update/{mType}/{mName}/{mValue}", SaveMetricHandler(db))
	r.Post("/update/", SaveJSONMetricHandler(db, config.HashKey))
	r.Post("/updates/", SaveBatchJSONMetricHandler(db, config.HashKey))
	return r
}
