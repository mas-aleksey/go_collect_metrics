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
}

type ServerConfig struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func LookupString(envName, defaultValue string) string {
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		return valueEnv
	} else {
		return defaultValue
	}
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

func MakeAgentConfig(address string, reportInterval time.Duration, pollInterval time.Duration) AgentConfig {
	cfg := AgentConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.ReportInterval = LookupDuration("REPORT_INTERVAL", reportInterval)
	cfg.PollInterval = LookupDuration("POLL_INTERVAL", pollInterval)
	return cfg
}

func MakeServerConfig(address string, restore bool, storeInterval time.Duration, storeFile string) ServerConfig {
	cfg := ServerConfig{}
	cfg.Address = LookupString("ADDRESS", address)
	cfg.Restore = LookupBool("RESTORE", restore)
	cfg.StoreInterval = LookupDuration("STORE_INTERVAL", storeInterval)
	cfg.StoreFile = LookupString("STORE_FILE", storeFile)
	return cfg
}
