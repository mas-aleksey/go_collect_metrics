package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		message    string
	}
	tests := []struct {
		name    string
		method  string
		request string
		db      storage.Storage
		want    want
	}{
		{
			name:    "check 405 not allowed",
			method:  http.MethodGet,
			request: "/update/type/name/value",
			db:      nil,
			want: want{
				statusCode: 405,
				message:    "",
			},
		},
		{
			name:    "check 404 wrong path",
			method:  http.MethodPost,
			request: "/update/type/name",
			db:      nil,
			want: want{
				statusCode: 404,
				message:    "404 page not found\n",
			},
		},
		{
			name:    "check 501 invalid metric type",
			method:  http.MethodPost,
			request: "/update/type/name/value",
			db:      nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:    "check 400 invalid gauge metric value",
			method:  http.MethodPost,
			request: "/update/gauge/Alloc/value",
			db:      nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:    "check 400 invalid counter metric value",
			method:  http.MethodPost,
			request: "/update/counter/PollCount/value",
			db:      nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:    "check 200 gauge success",
			method:  http.MethodPost,
			request: "/update/gauge/Alloc/123.456",
			db:      storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    "",
			},
		},
		{
			name:    "check 200 counter success",
			method:  http.MethodPost,
			request: "/update/counter/PollCount/123",
			db:      storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.db, utils.ServerConfig{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			request, _ := http.NewRequest(tt.method, ts.URL+tt.request, nil)
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
