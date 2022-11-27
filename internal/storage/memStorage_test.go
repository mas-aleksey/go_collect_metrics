package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"testing"
)

func TestMemStorage_SaveMetric(t *testing.T) {
	type want struct {
		gaugeMetrics   map[string]float64
		counterMetrics map[string]int64
	}
	tests := []struct {
		name    string
		metrics []utils.Metric
		want    want
	}{
		{
			name: "save metrics",
			metrics: []utils.Metric{
				utils.NewMetric("gauge", "RandomValue", "111.111"),
				utils.NewMetric("gauge", "RandomValue", "222.222"),
				utils.NewMetric("gauge", "RandomValue", "333.333"),
				utils.NewMetric("gauge", "Alloc", "123.456"),
				utils.NewMetric("gauge", "Frees", "1"),
				utils.NewMetric("gauge", "Frees", "0"),
				utils.NewMetric("gauge", "Sys", "555"),
				utils.NewMetric("counter", "PollCount", "1"),
				utils.NewMetric("counter", "PollCount", "2"),
				utils.NewMetric("counter", "PollCount", "3"),
			},
			want: want{
				gaugeMetrics:   map[string]float64{"RandomValue": 333.333, "Alloc": 123.456, "Frees": 0, "Sys": 555},
				counterMetrics: map[string]int64{"PollCount": 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemStorage()
			for _, metric := range tt.metrics {
				m.SaveMetric(metric)
			}
			assert.Equal(t, m.GaugeMetrics, tt.want.gaugeMetrics)
			assert.Equal(t, m.CounterMetrics, tt.want.counterMetrics)
		})
	}
}

func TestMemStorage_SaveJsonMetric(t *testing.T) {
	safeJSONMetric := func(body []byte) utils.JSONMetric {
		m, _ := utils.NewJsonMetric(body)
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
			m := NewMemStorage()
			for _, metric := range tt.metrics {
				m.SaveJSONMetric(metric)
			}
			assert.Equal(t, m.GaugeMetrics, tt.want.gaugeMetrics)
			assert.Equal(t, m.CounterMetrics, tt.want.counterMetrics)
		})
	}
}
