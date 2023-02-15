package handlers

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
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
	testStorage := storage.NewStorage(&utils.StorageConfig{})
	testStorage.GetBuffer().GaugeMetrics["Alloc"] = 123.456
	testStorage.GetBuffer().CounterMetrics["PoolCounter"] = 50

	tests := []struct {
		name     string
		method   string
		jsonData string
		hashKey  string
		db       storage.Storage
		want     want
	}{
		{
			name:     "check 405 not allowed",
			method:   http.MethodGet,
			jsonData: `{}`,
			hashKey:  "",
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
			hashKey:  "",
			db:       nil,
			want: want{
				statusCode: 422,
				message:    "invalid character 'H' looking for beginning of value\n",
			},
		},
		{
			name:     "check 501 invalid metric type",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"foo"}`,
			hashKey:  "",
			db:       nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:     "check 200 gauge success",
			method:   http.MethodPost,
			jsonData: `{"ID":"Alloc","type":"gauge"}`,
			hashKey:  "",
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"Alloc","type":"gauge","value":123.456}`,
			},
		},
		{
			name:     "check 200 counter success",
			method:   http.MethodPost,
			jsonData: `{"ID":"PoolCounter","type":"counter"}`,
			hashKey:  "",
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"PoolCounter","type":"counter","delta":50}`,
			},
		},
		{
			name:     "check 200 counter success with hash",
			method:   http.MethodPost,
			jsonData: `{"ID":"PoolCounter","type":"counter"}`,
			hashKey:  "some_key",
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `{"id":"PoolCounter","type":"counter","delta":50,"hash":"5a58361db229f5a67734904e9feb3e46e0183195fad9a9463f26a303790d2c2e"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.db, utils.ServerConfig{Address: "adr", HashKey: tt.hashKey})
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
	tests := []struct {
		name     string
		hashKey  string
		waitLen  string
		waitBody string
	}{
		{
			name:     "without hash key",
			hashKey:  "",
			waitLen:  "73",
			waitBody: `{"id":"Alloc","type":"gauge","value":123.456}`,
		},
		{
			name:     "with hash key",
			hashKey:  "some_key",
			waitLen:  "125",
			waitBody: `{"id":"Alloc","type":"gauge","value":123.456,"hash":"0dc3f697647c164dda7cb9fffc26aa833d4cb1fbec64a04f783cc57858d986d0"}`,
		},
	}
	testStorage := storage.NewStorage(&utils.StorageConfig{})
	testStorage.GetBuffer().GaugeMetrics["Alloc"] = 123.456
	testStorage.GetBuffer().CounterMetrics["PoolCounter"] = 50

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := GetRouter(testStorage, utils.ServerConfig{Address: "adr", HashKey: tt.hashKey})
			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			var body = []byte(`{"ID":"Alloc","type":"gauge"}`)
			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			_, err := w.Write(body)
			require.NoError(t, err)
			err = w.Close()
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, ts.URL+"/value/", &b)
			require.NoError(t, err)

			request.Header.Set("Accept-Encoding", "gzip")
			request.Header.Set("Content-Encoding", "gzip")
			request.Header.Set("Content-Type", "application/json")
			result, err := client.Do(request)
			require.NoError(t, err)

			assert.Equal(t, 200, result.StatusCode)
			assert.Equal(t, "gzip", result.Header.Get("Content-Encoding"))
			assert.Equal(t, tt.waitLen, result.Header.Get("Content-Length"))

			gzipReader, err := gzip.NewReader(result.Body)
			require.NoError(t, err)

			resBody, err := io.ReadAll(gzipReader)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)
			err = gzipReader.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.waitBody, string(resBody))
		})
	}
}
