package storage

import (
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

type Storage interface {
	Init() error
	Close()
	GetConfig() *utils.StorageConfig
	GetBuffer() *Buffer
	Ping() bool
	Save()
	SaveIfSyncMode()
}

func NewStorage(config *utils.StorageConfig) Storage {
	if config.DatabaseDSN != "" {
		return &PgStorage{
			Buffer: NewBuffer(),
			Config: config,
		}
	} else {
		return &MemStorage{
			Buffer: NewBuffer(),
			Config: config,
		}
	}
}
