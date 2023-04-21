package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	HashKey        string
}

type ServerConfig struct {
	Address string
	HashKey string
}

type StorageConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	DatabaseDSN   string
}

func LookupString(envName, defaultValue string) string {
	var result string
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		result = valueEnv
	} else {
		result = defaultValue
	}
	log.Printf("env %s = %s", envName, result)
	return result
}

func LookupDuration(envName string, defaultValue time.Duration) time.Duration {
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := time.ParseDuration(valueEnv)
		if err != nil {
			log.Printf("Filed to parse time.Duration env: %s: %s", envName, err)
			return defaultValue
		}
		return value
	} else {
		return defaultValue
	}
}

func LookupBool(envName string, defaultValue bool) bool {
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := strconv.ParseBool(valueEnv)
		if err != nil {
			log.Printf("Filed to parse bool env: %s: %s", envName, err)
			return defaultValue
		}
		return value
	} else {
		return defaultValue
	}
}

func MakeAgentConfig(address string, reportInterval time.Duration, pollInterval time.Duration, hashKey string) AgentConfig {
	cfg := AgentConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.ReportInterval = LookupDuration("REPORT_INTERVAL", reportInterval)
	cfg.PollInterval = LookupDuration("POLL_INTERVAL", pollInterval)
	cfg.HashKey = LookupString("KEY", hashKey)
	return cfg
}

func MakeServerConfig(address, hashKey string) ServerConfig {
	cfg := ServerConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.HashKey = LookupString("KEY", hashKey)
	return cfg
}

func MakeStorageConfig(restore bool, storeInterval time.Duration, storeFile, databaseDSN string) StorageConfig {
	cfg := StorageConfig{}
	cfg.Restore = LookupBool("RESTORE", restore)
	cfg.StoreInterval = LookupDuration("STORE_INTERVAL", storeInterval)
	cfg.StoreFile = LookupString("STORE_FILE", storeFile)
	cfg.DatabaseDSN = LookupString("DATABASE_DSN", databaseDSN)
	return cfg
}
