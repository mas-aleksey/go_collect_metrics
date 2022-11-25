package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

var RuntimeMetricNames = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc",
}

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
		fmt.Println("Unknown type", val, reflect.TypeOf(v))
		return "0"
	}
}
