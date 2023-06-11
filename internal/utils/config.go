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
	RateLimit      int
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
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		result = valueEnv
	}
	log.Printf("env %s = %s", envName, result)
	return result
}

func LookupInt(envName string, defaultValue int) int {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		intVar, err := strconv.Atoi(valueEnv)
		if err != nil {
			log.Printf("Filed to parse int env: %s: %s", envName, err)
		} else {
			result = intVar
		}
	}
	log.Printf("env %s = %d", envName, result)
	return result
}

func LookupDuration(envName string, defaultValue time.Duration) time.Duration {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := time.ParseDuration(valueEnv)
		if err != nil {
			log.Printf("Filed to parse time.Duration env: %s: %s", envName, err)
		} else {
			result = value
		}
	}
	log.Printf("env %s = %s", envName, result)
	return result
}

func LookupBool(envName string, defaultValue bool) bool {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := strconv.ParseBool(valueEnv)
		if err != nil {
			log.Printf("Filed to parse bool env: %s: %s", envName, err)
		} else {
			result = value
		}
	}
	return result
}

func MakeAgentConfig(address string, reportInterval time.Duration, pollInterval time.Duration, hashKey string, rateLimit int) AgentConfig {
	cfg := AgentConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.ReportInterval = LookupDuration("REPORT_INTERVAL", reportInterval)
	cfg.PollInterval = LookupDuration("POLL_INTERVAL", pollInterval)
	cfg.HashKey = LookupString("KEY", hashKey)
	cfg.RateLimit = LookupInt("RATE_LIMIT", rateLimit)
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
