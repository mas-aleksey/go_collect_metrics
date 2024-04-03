package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"
	"time"
)

// AgentConfig - структура конфигурации агента.
type AgentConfig struct {
	Address        string        `json:"address,omitempty"`
	ReportInterval time.Duration `json:"report_interval,omitempty"`
	PollInterval   time.Duration `json:"poll_interval,omitempty"`
	HashKey        string        `json:"hash_key,omitempty"`
	CryptoKey      string        `json:"crypto_key,omitempty"`
	RateLimit      int           `json:"rate_limit,omitempty"`
}

// ServerConfig - структура конфигурации сервера.
type ServerConfig struct {
	Address          string        `json:"address,omitempty"`
	HashKey          string        `json:"hash_key,omitempty"`
	CryptoKey        string        `json:"crypto_key,omitempty"`
	TrustedSubnet    string        `json:"trusted_subnet,omitempty"`
	TrustedNetPrefix *netip.Prefix `json:"-"`
}

// StorageConfig - структура конфигурации хранилища.
type StorageConfig struct {
	StoreInterval time.Duration `json:"store_interval,omitempty"`
	StoreFile     string        `json:"store_file,omitempty"`
	Restore       bool          `json:"restore,omitempty"`
	DatabaseDSN   string        `json:"database_dsn,omitempty"`
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

func loadFromFile(filePath string, v any) error {
	if filePath == "" {
		return nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func lookupString(flagName, envName, valueFromConfigFile, defaultFlagValue string) string {
	var result string
	if valueFromConfigFile != "" {
		result = valueFromConfigFile
	} else {
		result = defaultFlagValue
	}
	if isFlagPassed(flagName) {
		result = defaultFlagValue
	}
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		result = valueEnv
	}
	log.Printf("env %s = %v", envName, result)
	return result
}

func lookupInt(flagName, envName string, valueFromConfigFile, defaultFlagValue int) (int, error) {
	var result int
	if valueFromConfigFile != 0 {
		result = valueFromConfigFile
	} else {
		result = defaultFlagValue
	}
	if isFlagPassed(flagName) {
		result = defaultFlagValue
	}
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		intVar, err := strconv.Atoi(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "int", err)
		} else {
			result = intVar
		}
	}
	log.Printf("env %s = %v", envName, result)
	return result, nil
}

func lookupDuration(flagName, envName string, valueFromConfigFile, defaultFlagValue time.Duration) (time.Duration, error) {
	var result time.Duration
	if valueFromConfigFile != 0 {
		result = valueFromConfigFile
	} else {
		result = defaultFlagValue
	}
	if isFlagPassed(flagName) {
		result = defaultFlagValue
	}
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := time.ParseDuration(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "duration", err)
		} else {
			result = value
		}
	}
	log.Printf("env %s = %v", envName, result)
	return result, nil
}

func lookupBool(flagName, envName string, valueFromConfigFile, defaultFlagValue bool) (bool, error) {
	result := valueFromConfigFile
	if isFlagPassed(flagName) {
		result = defaultFlagValue
	}
	valueEnv, ok := os.LookupEnv(envName)
	if ok {
		value, err := strconv.ParseBool(valueEnv)
		if err != nil {
			return result, newEnvError(envName, "bool", err)
		} else {
			result = value
		}
	}
	log.Printf("env %s = %v", envName, result)
	return result, nil
}

// MakeAgentConfig - метод создания конфигурации агента.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeAgentConfig(
	configFile string, address string, reportInterval time.Duration,
	pollInterval time.Duration, hashKey string, cryptoKey string, rateLimit int,
) (AgentConfig, error) {

	var err error = nil
	cfg := AgentConfig{}
	configFile = lookupString("config", "CONFIG", "", configFile)
	err = loadFromFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Address = lookupString("a", "ADDRESS", cfg.Address, address)

	cfg.ReportInterval, err = lookupDuration("r", "REPORT_INTERVAL", cfg.ReportInterval, reportInterval)
	if err != nil {
		return cfg, err
	}
	cfg.PollInterval, err = lookupDuration("p", "POLL_INTERVAL", cfg.PollInterval, pollInterval)
	if err != nil {
		return cfg, err
	}
	cfg.HashKey = lookupString("k", "KEY", cfg.HashKey, hashKey)
	cfg.CryptoKey = lookupString("crypto-key", "CRYPTO_KEY", cfg.CryptoKey, cryptoKey)
	cfg.RateLimit, err = lookupInt("l", "RATE_LIMIT", cfg.RateLimit, rateLimit)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

// MakeServerConfig - метод создания конфигурации сервера.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeServerConfig(configFile, address, hashKey, cryptoKey, trustedSubnet string) (ServerConfig, error) {
	var err error = nil
	cfg := ServerConfig{}
	configFile = lookupString("config", "CONFIG", "", configFile)
	err = loadFromFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Address = lookupString("a", "ADDRESS", cfg.Address, address)
	cfg.HashKey = lookupString("k", "KEY", cfg.HashKey, hashKey)
	cfg.CryptoKey = lookupString("crypto-key", "CRYPTO_KEY", cfg.CryptoKey, cryptoKey)
	cfg.TrustedSubnet = lookupString("t", "TRUSTED_SUBNET", cfg.TrustedSubnet, trustedSubnet)
	if cfg.TrustedSubnet != "" {
		network, err := netip.ParsePrefix(cfg.TrustedSubnet)
		if err != nil {
			return cfg, err
		}
		cfg.TrustedNetPrefix = &network
	}
	return cfg, nil
}

// MakeStorageConfig - метод создания конфигурации хранилища.
// значения, переданные через параметры запуска, переопределяются значениями из переменных окружения.
func MakeStorageConfig(configFile string, restore bool, storeInterval time.Duration, storeFile, databaseDSN string) (StorageConfig, error) {
	var err error = nil
	cfg := StorageConfig{}
	configFile = lookupString("config", "CONFIG", "", configFile)
	err = loadFromFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Restore, err = lookupBool("r", "RESTORE", cfg.Restore, restore)
	if err != nil {
		return cfg, err
	}
	cfg.StoreInterval, err = lookupDuration("i", "STORE_INTERVAL", cfg.StoreInterval, storeInterval)
	if err != nil {
		return cfg, err
	}
	cfg.StoreFile = lookupString("f", "STORE_FILE", cfg.StoreFile, storeFile)
	cfg.DatabaseDSN = lookupString("d", "DATABASE_DSN", cfg.DatabaseDSN, databaseDSN)
	return cfg, nil
}
