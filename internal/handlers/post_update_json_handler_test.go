package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveJsonMetricHandler(t *testing.T) {
	testStorage := storage.NewStorage(&utils.StorageConfig{})
	testStorage.GetBuffer().CounterMetrics["PoolCounter"] = 77
	type want struct {
		statusCode int
		message    string
	}
	tests := []struct {
		name     string
		method   string
		jsonData string
		db       storage.Storage
		want     want
	}{
		{
			name:     "check 405 not allowed",
			method:   http.MethodGet,
			jsonData: `{}`,
			db:       nil,
			want: want{
				statusCode: 405,
				message:    "",
			},
		},
		{
			name:     "check 422 invalid json body",
			method:   http.MethodPost,
			jsonData: `{"message":Hello}`,
			db:       nil,
			want: want{
				statusCode: 422,
				message:    "invalid character 'H' looking for beginning of value\n",
			},
		},
		{
			name:     "check 501 invalid metric type",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"foo","Value":123}`,
			db:       nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:     "check 422 invalid gauge metric value",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"gauge","Value":"foo"}`,
			db:       nil,
			want: want{
				statusCode: 422,
				message:    "json: cannot unmarshal string into Go struct field JSONMetric.value of type float64\n",
			},
		},
		{
			name:     "check 400 nil gauge metric value",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"gauge"}`,
			db:       nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:     "check 422 invalid counter metric value",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"counter","Delta":"foo"}`,
			db:       nil,
			want: want{
				statusCode: 422,
				message:    "json: cannot unmarshal string into Go struct field JSONMetric.delta of type int64\n",
			},
		},
		{
			name:     "check 400 nil counter metric value",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"counter"}`,
			db:       nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:     "check 200 gauge success",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"gauge","Value":123.456}`,
			db:       storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    `{"id":"Alloc","type":"gauge","value":123.456,"hash":"9364e01ae0e8cf907ef330aa4a9691ad1e68aa51f05f4e2e73f23b406d0ce36a"}`,
			},
		},
		{
			name:     "check 200 counter success",
			method:   http.MethodPost,
			jsonData: `{"ID":"PoolCounter","type":"counter","Delta":123}`,
			db:       storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    `{"id":"PoolCounter","type":"counter","delta":123,"hash":"2799917354025ae1c468eb210efe049b9818c087b6ca186b87812e382952bdcf"}`,
			},
		},
		{
			name:     "check 200 counter success invalid hash",
			method:   http.MethodPost,
			jsonData: `{"id":"PoolCounter","type":"counter","delta":123,"hash":"some_text"}`,
			db:       storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 400,
				message:    "Invalid metric hash\n",
			},
		},
		{
			name:     "check 200 counter success valid hash",
			method:   http.MethodPost,
			jsonData: `{"id":"PoolCounter","type":"counter","delta":123,"hash":"2799917354025ae1c468eb210efe049b9818c087b6ca186b87812e382952bdcf"}`,
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"PoolCounter","type":"counter","delta":200,"hash":"e1265e8f3d1ecc83a870f5d5ee7a06c5b85393eed91d85e949dbf5bf4c44c765"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.db, utils.ServerConfig{Address: "adr", HashKey: "key"})
			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			var body = []byte(tt.jsonData)
			request, _ := http.NewRequest(tt.method, ts.URL+"/update/", bytes.NewBuffer(body))
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
