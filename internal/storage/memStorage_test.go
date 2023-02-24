package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"testing"
)

func TestMemStorage_SaveMetric(t *testing.T) {

	makeMetric := func(metricType, metricName, metricValue string) utils.JSONMetric {
		m, _ := utils.NewJSONMetric(metricType, metricName, metricValue)
		return m
	}
	type want struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name    string
		metrics []utils.JSONMetric
		want    want
	}{
		{
			name: "save metrics",
			metrics: []utils.JSONMetric{
				makeMetric("gauge", "RandomValue", "111.111"),
				makeMetric("gauge", "RandomValue", "222.222"),
				makeMetric("gauge", "RandomValue", "333.333"),
				makeMetric("gauge", "Alloc", "123.456"),
				makeMetric("gauge", "Frees", "1"),
				makeMetric("gauge", "Frees", "0"),
				makeMetric("gauge", "Sys", "555"),
				makeMetric("counter", "PollCount", "1"),
				makeMetric("counter", "PollCount", "2"),
				makeMetric("counter", "PollCount", "3"),
			},
			want: want{
				gaugeMetrics:   map[string]float64{"RandomValue": 333.333, "Alloc": 123.456, "Frees": 0, "Sys": 555},
				counterMetrics: map[string]int64{"PollCount": 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewStorage(&utils.StorageConfig{})
			for _, metric := range tt.metrics {
				db.GetBuffer().PutJSONMetric(metric)
			}
			assert.Equal(t, db.GetBuffer().GaugeMetrics, tt.want.gaugeMetrics)
			assert.Equal(t, db.GetBuffer().CounterMetrics, tt.want.counterMetrics)
		})
	}
}

func TestMemStorage_SaveJsonMetric(t *testing.T) {
	safeJSONMetric := func(body []byte) utils.JSONMetric {
		m, _ := utils.LoadJSONMetric(body)
		return m
	}
	type want struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name    string
		metrics []utils.JSONMetric
		want    want
	}{
		{
			name: "save json metrics",
			metrics: []utils.JSONMetric{
				safeJSONMetric([]byte(`{"ID":"RandomValue","type":"gauge","Value":111.111}`)),
				safeJSONMetric([]byte(`{"ID":"RandomValue","type":"gauge","Value":222.222}`)),
				safeJSONMetric([]byte(`{"ID":"RandomValue","type":"gauge","Value":333.333}`)),
				safeJSONMetric([]byte(`{"ID":"Alloc","type":"gauge","Value":123.456}`)),
				safeJSONMetric([]byte(`{"ID":"Frees","type":"gauge","Value":1}`)),
				safeJSONMetric([]byte(`{"ID":"Frees","type":"gauge","Value":0}`)),
				safeJSONMetric([]byte(`{"ID":"Sys","type":"gauge","Value":555}`)),
				safeJSONMetric([]byte(`{"ID":"PollCount","type":"counter","Delta":1}`)),
				safeJSONMetric([]byte(`{"ID":"PollCount","type":"counter","Delta":2}`)),
				safeJSONMetric([]byte(`{"ID":"PollCount","type":"counter","Delta":3}`)),
			},
			want: want{
				gaugeMetrics:   map[string]float64{"RandomValue": 333.333, "Alloc": 123.456, "Frees": 0, "Sys": 555},
				counterMetrics: map[string]int64{"PollCount": 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewStorage(&utils.StorageConfig{})
			for _, metric := range tt.metrics {
				db.GetBuffer().PutJSONMetric(metric)
			}
			assert.Equal(t, db.GetBuffer().GaugeMetrics, tt.want.gaugeMetrics)
			assert.Equal(t, db.GetBuffer().CounterMetrics, tt.want.counterMetrics)
		})
	}
}

func TestMemStorage_SetMetricValue(t *testing.T) {
	m := MemStorage{
		Buffer: &Buffer{
			GaugeMetrics:   map[string]float64{"name": 123.4},
			CounterMetrics: map[string]int64{"name": 123},
		},
	}
	gaugeMetric := utils.JSONMetric{
		ID:    "name",
		MType: "gauge",
	}
	counterMetric := utils.JSONMetric{
		ID:    "name",
		MType: "counter",
	}
	m.GetBuffer().UpdateJSONMetricValue(&gaugeMetric)
	m.GetBuffer().UpdateJSONMetricValue(&counterMetric)
	assert.Equal(t, 123.4, *gaugeMetric.Value)
	assert.Equal(t, int64(123), *counterMetric.Delta)
}

func TestMemStorage_SetJSONMetricValue(t *testing.T) {
	m := MemStorage{
		Buffer: &Buffer{
			GaugeMetrics:   map[string]float64{"name": 123.4},
			CounterMetrics: map[string]int64{"name": 123},
		},
	}
	gaugeMetric := utils.JSONMetric{
		ID:    "name",
		MType: "gauge",
	}
	counterMetric := utils.JSONMetric{
		ID:    "name",
		MType: "counter",
	}
	m.GetBuffer().UpdateJSONMetricValue(&counterMetric)
	m.GetBuffer().UpdateJSONMetricValue(&gaugeMetric)
	assert.Equal(t, int64(123), *counterMetric.Delta)
	assert.Equal(t, 123.4, *gaugeMetric.Value)
}
