package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveBatchJSONMetricHandler(t *testing.T) {
	testStorage := storage.NewStorage(&utils.StorageConfig{})
	testStorage.GetBuffer().CounterMetrics["PoolCounter"] = 77
	type want struct {
		statusCode int
		message    string
	}
	tests := []struct {
		name     string
		jsonData string
		db       storage.Storage
		want     want
	}{
		{
			name:     "check 200 batch metrics come",
			jsonData: `[{"id":"PoolCounter","type":"counter","delta":123,"hash":"2799917354025ae1c468eb210efe049b9818c087b6ca186b87812e382952bdcf"}]`,
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `[{"id":"PoolCounter","type":"counter","delta":200,"hash":"e1265e8f3d1ecc83a870f5d5ee7a06c5b85393eed91d85e949dbf5bf4c44c765"}]`,
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
			request, _ := http.NewRequest(http.MethodPost, ts.URL+"/updates/", bytes.NewBuffer(body))
			result, err := client.Do(request)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.message, string(resBody))

			var metrics []utils.JSONMetric
			err = json.Unmarshal(resBody, &metrics)
			assert.Nil(t, err)
		})
	}
}
