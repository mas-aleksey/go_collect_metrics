package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetValueJSONMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		message    string
	}
	testStorage := storage.NewMemStorage()
	testStorage.GaugeMetrics["Alloc"] = 123.456
	testStorage.CounterMetrics["PoolCounter"] = 50

	tests := []struct {
		name       string
		method     string
		jsonData   string
		memStorage *storage.MemStorage
		want       want
	}{
		{
			name:       "check 405 not allowed",
			method:     http.MethodGet,
			jsonData:   `{}`,
			memStorage: nil,
			want: want{
				statusCode: 405,
				message:    "",
			},
		},
		{
			name:       "check 422 invalid json body",
			method:     http.MethodPost,
			jsonData:   `{"message":Hello}`,
			memStorage: nil,
			want: want{
				statusCode: 422,
				message:    "invalid character 'H' looking for beginning of value\n",
			},
		},
		{
			name:       "check 501 invalid metric type",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"foo"}`,
			memStorage: nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:       "check 200 gauge success",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"gauge"}`,
			memStorage: testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"Alloc","type":"gauge","value":123.456}`,
			},
		},
		{
			name:       "check 200 counter success",
			method:     http.MethodPost,
			jsonData:   `{"ID":"PoolCounter","type":"counter"}`,
			memStorage: testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"PoolCounter","type":"counter","delta":50}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.memStorage)
			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			var body = []byte(tt.jsonData)
			request, err := http.NewRequest(tt.method, ts.URL+"/value/", bytes.NewBuffer(body))
			require.NoError(t, err)

			result, err := client.Do(request)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}

func TestCompressedSetValueJSONMetricHandler(t *testing.T) {
	testStorage := storage.NewMemStorage()
	testStorage.GaugeMetrics["Alloc"] = 123.456
	testStorage.CounterMetrics["PoolCounter"] = 50

	r := GetRouter(testStorage)
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{}

	var body = []byte(`{"ID":"Alloc","type":"gauge"}`)
	request, err := http.NewRequest(http.MethodPost, ts.URL+"/value/", bytes.NewBuffer(body))
	require.NoError(t, err)

	request.Header.Set("Accept-Encoding", "gzip")
	result, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "gzip", result.Header.Get("Content-Encoding"))
	assert.Equal(t, "73", result.Header.Get("Content-Length"))

	resBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	err = result.Body.Close()
	require.NoError(t, err)

	assert.NotNil(t, resBody)
}
