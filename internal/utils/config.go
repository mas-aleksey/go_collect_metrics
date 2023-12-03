package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// AgentConfig - структура конфигурации агента.
type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	HashKey        string
	CryptoKey      string
	RateLimit      int
}

// ServerConfig - структура конфигурации сервера.
type ServerConfig struct {
	Address   string
	HashKey   string
	CryptoKey string
}

// StorageConfig - структура конфигурации хранилища.
type StorageConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	DatabaseDSN   string
}

// EnvError - тип ошибки, связанный с получением переменной из окружения.
type EnvError struct {
	EnvName string
	EnvType string
	Err     error
}

// метод преобразования ошибки EnvError к строке.
func (e *EnvError) Error() string {
	return fmt.Sprintf("Filed to parse %s env: %s: %v", e.EnvType, e.EnvName, e.Err)
}

func newEnvError(eName, eType string, err error) error {
	return &EnvError{
		EnvName: eName,
		EnvType: eType,
		Err:     err,
	}
}

func lookupString(envName, defaultValue string) string {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		result = valueEnv
	}
	log.Printf("env %s = %s", envName, result)
	return result
}

func lookupInt(envName string, defaultValue int) (int, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		intVar, err := strconv.Atoi(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "int", err)
		} else {
			result = intVar
		}
	}
	log.Printf("env %s = %d", envName, result)
	return result, nil
}

func lookupDuration(envName string, defaultValue time.Duration) (time.Duration, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := time.ParseDuration(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "int", err)
		} else {
			result = value
		}
	}
	log.Printf("env %s = %s", envName, result)
	return result, nil
}

func lookupBool(envName string, defaultValue bool) (bool, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := strconv.ParseBool(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "int", err)
		} else {
			result = value
		}
	}
	return result, nil
}

// MakeAgentConfig - метод создания конфигурации агента.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeAgentConfig(address string, reportInterval time.Duration, pollInterval time.Duration, hashKey string, cryptoKey string, rateLimit int) (AgentConfig, error) {
	var err error = nil
	cfg := AgentConfig{}
	cfg.Address = lookupString("ADDRESS", address)
	cfg.ReportInterval, err = lookupDuration("REPORT_INTERVAL", reportInterval)
	if err != nil {
		return cfg, err
	}
	cfg.PollInterval, err = lookupDuration("POLL_INTERVAL", pollInterval)
	if err != nil {
		return cfg, err
	}
	cfg.HashKey = lookupString("KEY", hashKey)
	cfg.CryptoKey = lookupString("CRYPTO_KEY", cryptoKey)
	cfg.RateLimit, err = lookupInt("RATE_LIMIT", rateLimit)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// MakeServerConfig - метод создания конфигурации сервера.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeServerConfig(address, hashKey, cryptoKey string) ServerConfig {
	cfg := ServerConfig{}
	cfg.Address = lookupString("ADDRESS", address)
	cfg.HashKey = lookupString("KEY", hashKey)
	cfg.CryptoKey = lookupString("CRYPTO_KEY", cryptoKey)
	return cfg
}

// MakeStorageConfig - метод создания конфигурации хранилища.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeStorageConfig(restore bool, storeInterval time.Duration, storeFile, databaseDSN string) (StorageConfig, error) {
	var err error = nil
	cfg := StorageConfig{}
	cfg.Restore, err = lookupBool("RESTORE", restore)
	if err != nil {
		return cfg, err
	}
	cfg.StoreInterval, err = lookupDuration("STORE_INTERVAL", storeInterval)
	if err != nil {
		return cfg, err
	}
	cfg.StoreFile = lookupString("STORE_FILE", storeFile)
	cfg.DatabaseDSN = lookupString("DATABASE_DSN", databaseDSN)
	return cfg, nil
}
