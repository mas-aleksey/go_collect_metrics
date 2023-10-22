package utils

import (
	"fmt"
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

type EnvError struct {
	EnvName string
	EnvType string
	Err     error
}

func (e *EnvError) Error() string {
	return fmt.Sprintf("Filed to parse %s env: %s: %v", e.EnvType, e.EnvName, e.Err)
}

func NewEnvError(eName, eType string, err error) error {
	return &EnvError{
		EnvName: eName,
		EnvType: eType,
		Err:     err,
	}
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

func LookupInt(envName string, defaultValue int) (int, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		intVar, err := strconv.Atoi(valueEnv)
		if err != nil {
			return result, NewEnvError(envName, "int", err)
		} else {
			result = intVar
		}
	}
	log.Printf("env %s = %d", envName, result)
	return result, nil
}

func LookupDuration(envName string, defaultValue time.Duration) (time.Duration, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := time.ParseDuration(valueEnv)
		if err != nil {
			return result, NewEnvError(envName, "int", err)
		} else {
			result = value
		}
	}
	log.Printf("env %s = %s", envName, result)
	return result, nil
}

func LookupBool(envName string, defaultValue bool) (bool, error) {
	result := defaultValue
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := strconv.ParseBool(valueEnv)
		if err != nil {
			return result, NewEnvError(envName, "int", err)
		} else {
			result = value
		}
	}
	return result, nil
}

func MakeAgentConfig(address string, reportInterval time.Duration, pollInterval time.Duration, hashKey string, rateLimit int) (AgentConfig, error) {
	var err error = nil
	cfg := AgentConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.ReportInterval, err = LookupDuration("REPORT_INTERVAL", reportInterval)
	if err != nil {
		return cfg, err
	}
	cfg.PollInterval, err = LookupDuration("POLL_INTERVAL", pollInterval)
	if err != nil {
		return cfg, err
	}
	cfg.HashKey = LookupString("KEY", hashKey)
	cfg.RateLimit, err = LookupInt("RATE_LIMIT", rateLimit)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func MakeServerConfig(address, hashKey string) ServerConfig {
	cfg := ServerConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.HashKey = LookupString("KEY", hashKey)
	return cfg
}

func MakeStorageConfig(restore bool, storeInterval time.Duration, storeFile, databaseDSN string) (StorageConfig, error) {
	var err error = nil
	cfg := StorageConfig{}
	cfg.Restore, err = LookupBool("RESTORE", restore)
	if err != nil {
		return cfg, err
	}
	cfg.StoreInterval, err = LookupDuration("STORE_INTERVAL", storeInterval)
	if err != nil {
		return cfg, err
	}
	cfg.StoreFile = LookupString("STORE_FILE", storeFile)
	cfg.DatabaseDSN = LookupString("DATABASE_DSN", databaseDSN)
	return cfg, nil
}
