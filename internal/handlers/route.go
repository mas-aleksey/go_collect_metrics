// Package handlers - функционал обработчиков API методов сервера
package handlers

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// GetRouter - метод регистрирует роуты для сервера.
func GetRouter(db storage.Storage, config utils.ServerConfig, privateKey *utils.PrivateKey) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(1, "application/json", "text/html", "text/plain"))
	r.Use(middleware.AllowContentEncoding("gzip"))
	if config.TrustedNetPrefix != nil {
		r.Use(middleware.RealIP)
		r.Use(CheckTrustedSubnet(config.TrustedNetPrefix))
	}
	r.Get("/", IndexHandler(db))
	r.Get("/ping", GetPingHandler(db))
	r.Get("/value/{mType}/{mName}", GetValueMetricHandler(db))
	r.Post("/update/{mType}/{mName}/{mValue}", SaveMetricHandler(db))
	r.Post("/value/", GetJSONMetricHandler(db, config.HashKey, privateKey))
	r.Post("/update/", SaveJSONMetricHandler(db, config.HashKey, privateKey))
	r.Post("/updates/", SaveBatchJSONMetricHandler(db, config.HashKey, privateKey))
	return r
}
