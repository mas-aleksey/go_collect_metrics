package clients

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetricClient(t *testing.T) {
	expected := MetricClient{baseURL: "localhost:8080", client: &http.Client{}}
	result := NewMetricClient("localhost:8080")
	assert.Equal(t, result, expected)
}

func TestMetricClient_postMetric(t *testing.T) {
	metric := utils.NewMetric("type", "name", "value")
	wantURL := "/update/type/name/value"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, wantURL)
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.Header.Get("Content-Type"), "text/plain")
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	mc := NewMetricClient(svr.URL)
	out, err := mc.postMetric(metric)
	assert.Nil(t, err)
	assert.Equal(t, out, "")
}

func TestMetricClient_SendMetrics(t *testing.T) {
	expectedCount := 29
	statistic := utils.NewStatistic()
	urls := make([]string, 0)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.Path)
		urls = append(urls, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	mc := NewMetricClient(svr.URL)
	mc.SendMetrics(*statistic)
	assert.Equal(t, len(urls), expectedCount)
}
