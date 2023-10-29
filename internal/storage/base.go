// Package storage - реализация взаимодействие с объектом Storage.
package storage

import (
	"context"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

// Storage - общий интерфейс для взаимодействия с любым типом хранилища.
type Storage interface {
	// Init инициализация подключения
	Init(context.Context) error
	// Close закрытие соединения
	Close(context.Context)
	// Ping проверка открытого соединения
	Ping(context.Context) bool
	// UpdateJSONMetric обновление одной метрики
	UpdateJSONMetric(context.Context, utils.JSONMetric) (utils.JSONMetric, error)
	// UpdateJSONMetrics обновление списка метрик
	UpdateJSONMetrics(context.Context, []utils.JSONMetric) ([]utils.JSONMetric, error)
	// GetJSONMetric получение одной метрики
	GetJSONMetric(context.Context, string, string) (utils.JSONMetric, error)
	// GetAllMetrics получение всех метрик
	GetAllMetrics(context.Context) ([]utils.JSONMetric, error)
}

// NewStorage - метод для создания объекта Storage
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
