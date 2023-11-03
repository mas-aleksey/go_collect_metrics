// Package utils - пакет для вспомогательных функций
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// RuntimeMetricNames - список имен для Runtime метрик
var RuntimeMetricNames = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc",
}

// ToStr - метод для приведения чисел в строку
func ToStr(v interface{}) string {

	switch val := v.(type) {
	case uint64:
		return strconv.FormatUint(val, 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', 3, 64)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	default:
		log.Println("Unknown type", val, reflect.TypeOf(v))
		return "0"
	}
}

// ToFloat64 - метод для приведения чисел к формату float64
func ToFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case uint64:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case uint32:
		return float64(val)
	default:
		log.Println("Unknown type", val, reflect.TypeOf(v))
		return float64(0)
	}
}

// CalcHash - метод для вычисления хеш-суммы
func CalcHash(data, hashKey string) *string {
	if hashKey == "" {
		return nil
	}
	h := hmac.New(sha256.New, []byte(hashKey))
	h.Write([]byte(data))
	dst := fmt.Sprintf("%x", h.Sum(nil))
	return &dst
}
