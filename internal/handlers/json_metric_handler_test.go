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

func TestSaveJsonMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		message    string
	}
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
			jsonData:   `{"ID":"Alloc","type":"foo","Value":123}`,
			memStorage: nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:       "check 422 invalid gauge metric value",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"gauge","Value":"foo"}`,
			memStorage: nil,
			want: want{
				statusCode: 422,
				message:    "json: cannot unmarshal string into Go struct field JsonMetric.value of type float64\n",
			},
		},
		{
			name:       "check 400 nil gauge metric value",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"gauge"}`,
			memStorage: nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:       "check 422 invalid counter metric value",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"counter","Delta":"foo"}`,
			memStorage: nil,
			want: want{
				statusCode: 422,
				message:    "json: cannot unmarshal string into Go struct field JsonMetric.delta of type int64\n",
			},
		},
		{
			name:       "check 400 nil counter metric value",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"counter"}`,
			memStorage: nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:       "check 200 gauge success",
			method:     http.MethodPost,
			jsonData:   `{"ID":"Alloc","type":"gauge","Value":123.456}`,
			memStorage: storage.NewMemStorage(),
			want: want{
				statusCode: 200,
				message:    "",
			},
		},
		{
			name:       "check 200 counter success",
			method:     http.MethodPost,
			jsonData:   `{"ID":"PoolCounter","type":"counter","Delta":123}`,
			memStorage: storage.NewMemStorage(),
			want: want{
				statusCode: 200,
				message:    "",
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
			request, _ := http.NewRequest(tt.method, ts.URL+"/update", bytes.NewBuffer(body))
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
