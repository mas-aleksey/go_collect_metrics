package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"storage"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		message     string
	}
	tests := []struct {
		name        string
		method      string
		contentType string
		request     string
		memStorage  *storage.MemStorage
		want        want
	}{
		{
			name:        "check 405 not allowed",
			method:      http.MethodGet,
			contentType: "text/plain",
			request:     "/update/type/name/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  405,
				message:     "Method GET not allowed\n",
			},
		},
		{
			name:        "check 404 wrong path",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/type/name",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
				message:     "Wrong path. Expected: 'HOST/update/metric_type/metric_name/metric_value'\n",
			},
		},
		{
			name:        "check 400 invalid ContentType",
			method:      http.MethodPost,
			contentType: "application/json",
			request:     "/update/type/name/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				message:     "Expected 'Content-Type' only 'text/plain'\n",
			},
		},
		{
			name:        "check 400 invalid metric type",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/type/name/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				message:     "Invalid metric type\n",
			},
		},
		{
			name:        "check 400 invalid metric name",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/gauge/name/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				message:     "Invalid metric name\n",
			},
		},
		{
			name:        "check 400 invalid gauge metric value",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/gauge/Alloc/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				message:     "Invalid metric value\n",
			},
		},
		{
			name:        "check 400 invalid counter metric value",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/counter/PollCount/value",
			memStorage:  nil,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				message:     "Invalid metric value\n",
			},
		},
		{
			name:        "check 200 gauge success",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/gauge/Alloc/123.456",
			memStorage:  storage.NewMemStorage(),
			want: want{
				contentType: "",
				statusCode:  200,
				message:     "",
			},
		},
		{
			name:        "check 200 counter success",
			method:      http.MethodPost,
			contentType: "text/plain",
			request:     "/update/counter/PollCount/123",
			memStorage:  storage.NewMemStorage(),
			want: want{
				contentType: "",
				statusCode:  200,
				message:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h := UpdateMetricHandler(tt.memStorage)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}
