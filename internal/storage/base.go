package storage

import (
	"context"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

type Storage interface {
	Init(context.Context) error
	Close(context.Context)
	Ping(context.Context) bool
	UpdateJSONMetric(context.Context, utils.JSONMetric) (utils.JSONMetric, error)
	UpdateJSONMetrics(context.Context, []utils.JSONMetric) ([]utils.JSONMetric, error)
	GetJSONMetric(context.Context, string, string) (utils.JSONMetric, error)
	GetAllMetrics(context.Context) ([]utils.JSONMetric, error)
}

func NewStorage(config *utils.StorageConfig) Storage {
	if config.DatabaseDSN != "" {
		return &PgStorage{
			Config: config,
		}
	} else {
		return &MemStorage{
			GaugeMetrics:   make(map[string]float64),
			CounterMetrics: make(map[string]int64),
			Config:         config,
		}
	}
}
