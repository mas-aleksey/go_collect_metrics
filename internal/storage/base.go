package storage

import (
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

type Storage interface {
	Init() error
	Close()
	Ping() bool
	UpdateJSONMetric(utils.JSONMetric) (utils.JSONMetric, error)
	UpdateJSONMetrics([]utils.JSONMetric) ([]utils.JSONMetric, error)
	GetJSONMetric(string, string) (utils.JSONMetric, error)
	GetAllMetrics() ([]utils.JSONMetric, error)
}

func NewStorage(config *utils.StorageConfig) Storage {
	config.DatabaseDSN = "postgresql://ml_platform_orchestrator_admin:pwd@localhost:5467/yandex"
	if config.DatabaseDSN != "" {
		db := &PgStorage{
			Config: config,
		}
		db.Init()
		return db
	} else {
		db := &MemStorage{
			GaugeMetrics:   make(map[string]float64),
			CounterMetrics: make(map[string]int64),
			Config:         config,
		}
		db.Init()
		return db
	}
}
