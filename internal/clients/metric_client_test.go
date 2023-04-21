package clients

import (
	"compress/gzip"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetricClient(t *testing.T) {
	expected := MetricClient{BaseClient: NewBaseClient("localhost:8080", 1*time.Second)}
	result := NewMetricClient("localhost:8080", 1*time.Second)
	assert.Equal(t, *result, expected)
}

func TestMetricClient_postMetric(t *testing.T) {
	tests := []struct {
		name     string
		metric   utils.JSONMetric
		hashKey  string
		wantBody string
	}{
		{
			name:     "empty counter hash",
			metric:   utils.NewCounterJSONMetric("name", 10),
			hashKey:  "",
			wantBody: "{\"id\":\"name\",\"type\":\"counter\",\"delta\":10}",
		},
		{
			name:     "fill counter hash",
			metric:   utils.NewCounterJSONMetric("name", 10),
			hashKey:  "123",
			wantBody: "{\"id\":\"name\",\"type\":\"counter\",\"delta\":10,\"hash\":\"775bc6d6bc40cb6535f865f85936c49453d224c5dc6047215247d17c28b6c8a1\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resBody, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				err = r.Body.Close()
				require.NoError(t, err)
				assert.Equal(t, tt.wantBody, string(resBody))
				assert.Equal(t, r.URL.Path, "/update/")
				assert.Equal(t, r.Method, "POST")
				assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
				w.WriteHeader(http.StatusOK)
			}))
			defer svr.Close()
			mc := NewMetricClient(svr.URL, 1*time.Second)
			tt.metric.Hash = utils.CalcHash(tt.metric.String(), tt.hashKey)
			body, err := json.Marshal(tt.metric)
			assert.Nil(t, err)
			err = mc.postBody(body, "update/", false)
			assert.Nil(t, err)
		})
	}
}

func TestMetricClient_SendJSONReport(t *testing.T) {
	statistic := utils.NewStatistic()
	report := utils.NewJSONReport(statistic, "")
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.Path)
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var metric utils.JSONMetric
		err = json.Unmarshal(body, &metric)
		assert.Nil(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	mc := NewMetricClient(svr.URL, 1*time.Second)
	err := mc.SendJSONReport(report)
	assert.Nil(t, err)
}

func TestMetricClient_SendBatchJSONReport(t *testing.T) {
	statistic := utils.NewStatistic()
	report := utils.NewJSONReport(statistic, "")
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/updates/", r.URL.Path)
		reader, err := gzip.NewReader(r.Body)
		assert.Nil(t, err)
		body, err := io.ReadAll(reader)
		assert.Nil(t, err)

		err = r.Body.Close()
		assert.Nil(t, err)
		err = reader.Close()
		assert.Nil(t, err)

		var metrics []utils.JSONMetric
		err = json.Unmarshal(body, &metrics)
		assert.Nil(t, err)

		assert.Equal(t, 29, len(metrics))
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	mc := NewMetricClient(svr.URL, 1*time.Second)
	err := mc.SendBatchJSONReport(report)
	assert.Nil(t, err)
}
