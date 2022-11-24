package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		message    string
	}
	tests := []struct {
		name       string
		method     string
		request    string
		memStorage *storage.MemStorage
		want       want
	}{
		{
			name:       "check 405 not allowed",
			method:     http.MethodGet,
			request:    "/update/type/name/value",
			memStorage: nil,
			want: want{
				statusCode: 405,
				message:    "Method GET not allowed\n",
			},
		},
		{
			name:       "check 404 wrong path",
			method:     http.MethodPost,
			request:    "/update/type/name",
			memStorage: nil,
			want: want{
				statusCode: 404,
				message:    "Wrong path. Expected: 'HOST/update/metric_type/metric_name/metric_value'\n",
			},
		},
		{
			name:       "check 502 invalid metric type",
			method:     http.MethodPost,
			request:    "/update/type/name/value",
			memStorage: nil,
			want: want{
				statusCode: 501,
				message:    "Invalid metric type\n",
			},
		},
		{
			name:       "check 400 invalid gauge metric value",
			method:     http.MethodPost,
			request:    "/update/gauge/Alloc/value",
			memStorage: nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:       "check 400 invalid counter metric value",
			method:     http.MethodPost,
			request:    "/update/counter/PollCount/value",
			memStorage: nil,
			want: want{
				statusCode: 400,
				message:    "Invalid metric value\n",
			},
		},
		{
			name:       "check 200 gauge success",
			method:     http.MethodPost,
			request:    "/update/gauge/Alloc/123.456",
			memStorage: storage.NewMemStorage(),
			want: want{
				statusCode: 200,
				message:    "",
			},
		},
		{
			name:       "check 200 counter success",
			method:     http.MethodPost,
			request:    "/update/counter/PollCount/123",
			memStorage: storage.NewMemStorage(),
			want: want{
				statusCode: 200,
				message:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			h := SaveMetricHandler(tt.memStorage)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}
