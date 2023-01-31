package storage

import (
	"encoding/json"
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"os"
)

type MemStorage struct {
	Buffer *Buffer
	Config *utils.StorageConfig
}

func (m *MemStorage) Init() error {
	if !m.Config.Restore {
		return fmt.Errorf("no need restore")
	}
	if m.Config.StoreFile == "" {
		return fmt.Errorf("filename is empty")
	}
	m.Buffer.Mutex.Lock()
	defer m.Buffer.Mutex.Unlock()

	data, err := os.ReadFile(m.Config.StoreFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, m.Buffer)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) Close() {
	m.Save()
}

func (m *MemStorage) GetConfig() *utils.StorageConfig {
	return m.Config
}

func (m *MemStorage) GetBuffer() *Buffer {
	return m.Buffer
}

func (m *MemStorage) Ping() bool {
	return true
}

func (m *MemStorage) Save() {
	err := m.saveToFile()
	if err != nil {
		log.Print("Failed save to file", err)
	} else {
		log.Print("Save storage to file")
	}
}

func (m *MemStorage) SaveIfSyncMode() {
	if m.Config.StoreInterval == 0 {
		m.Save()
	}
}

func (m *MemStorage) saveToFile() error {
	if m.Config.StoreFile == "" {
		return fmt.Errorf("filename is empty")
	}
	m.Buffer.Mutex.RLock()
	defer m.Buffer.Mutex.RUnlock()

	file, err := os.Create(m.Config.StoreFile)

	if err != nil {
		return err
	}
	data, err := json.Marshal(m.Buffer)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
