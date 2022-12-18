package clients

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetricClient(t *testing.T) {
	config := NewClientConfig("localhost:8080", 1*time.Second)
	expected := MetricClient{BaseClient: NewBaseClient(config)}
	result := NewMetricClient(config)
	assert.Equal(t, *result, expected)
}

func TestMetricClient_postMetric(t *testing.T) {
	metric := utils.NewCounterJSONMetric("name", 10)
	wantURL := "/update/"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, wantURL)
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	config := NewClientConfig(svr.URL, 1*time.Second)
	mc := NewMetricClient(config)
	err := mc.postJSONMetric(metric)
	assert.Nil(t, err)
}

func TestMetricClient_SendMetrics(t *testing.T) {
	statistic := utils.NewStatistic()
	report := utils.NewJSONReport(statistic)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	config := NewClientConfig(svr.URL, 1*time.Second)
	mc := NewMetricClient(config)
	err := mc.SendJSONReport(report)
	assert.Nil(t, err)
}
