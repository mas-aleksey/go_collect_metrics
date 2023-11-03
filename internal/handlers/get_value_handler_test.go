package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

func TestGetValueMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		message    string
	}
	testStorage := storage.NewStorage(&utils.StorageConfig{})
	_, _ = testStorage.UpdateJSONMetrics(context.Background(), []utils.JSONMetric{
		utils.NewGaugeJSONMetric("Alloc", 111.222),
		utils.NewCounterJSONMetric("PollCount", 333),
	})

	tests := []struct {
		name    string
		method  string
		request string
		db      storage.Storage
		want    want
	}{
		{
			name:    "check 405 not allowed",
			method:  http.MethodPost,
			request: "/value/type/name",
			db:      nil,
			want: want{
				statusCode: 405,
				message:    "",
			},
		},
		{
			name:    "check 404 wrong path",
			method:  http.MethodGet,
			request: "/value/type/name/foo",
			db:      nil,
			want: want{
				statusCode: 404,
				message:    "404 page not found\n",
			},
		},
		{
			name:    "check 404 gauge metric not found",
			method:  http.MethodGet,
			request: "/value/gauge/fooName",
			db:      testStorage,
			want: want{
				statusCode: 404,
				message:    "Metric not found\n",
			},
		},
		{
			name:    "check 404 counter metric not found",
			method:  http.MethodGet,
			request: "/value/counter/fooName",
			db:      testStorage,
			want: want{
				statusCode: 404,
				message:    "Metric not found\n",
			},
		},
		{
			name:    "check 200 gauge success",
			method:  http.MethodGet,
			request: "/value/gauge/Alloc",
			db:      testStorage,
			want: want{
				statusCode: 200,
				message:    "111.222",
			},
		},
		{
			name:    "check 200 counter success",
			method:  http.MethodGet,
			request: "/value/counter/PollCount",
			db:      testStorage,
			want: want{
				statusCode: 200,
				message:    "333",
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
