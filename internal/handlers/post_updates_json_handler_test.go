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

var DATA = `[
    {
        "id": "PollCount",
        "type": "counter",
        "delta": 10,
        "hash": "5848f030886e828a0481a33dd7386ee4c617c17b4f315645b83dd4b129514639"
    },
    {
        "id": "SomeParam",
        "type": "gauge",
        "value": 10.1,
        "hash": "4b37bfcaa5c419ce160ed8a1247c66a109b086336ac4d7d9df659a86c764a93b"
    },
    {
        "id": "PollCount",
        "type": "counter",
        "delta": 10,
        "hash": "5848f030886e828a0481a33dd7386ee4c617c17b4f315645b83dd4b129514639"
    },
    {
        "id": "SomeParam",
        "type": "gauge",
        "value": 11.1,
        "hash": "335bd95edd095efeb1d61bda8bc51f613690b8a6908866a2169a751bcf960933"
    }
]`

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
		hashKey  string
		db       storage.Storage
		want     want
	}{
		{
			name:     "check 200 batch metrics come",
			jsonData: `[{"id":"PoolCounter","type":"counter","delta":123,"hash":"2799917354025ae1c468eb210efe049b9818c087b6ca186b87812e382952bdcf"}]`,
			hashKey:  "key",
			db:       testStorage,
			want: want{
				statusCode: 200,
				message:    `[{"id":"PoolCounter","type":"counter","delta":200,"hash":"e1265e8f3d1ecc83a870f5d5ee7a06c5b85393eed91d85e949dbf5bf4c44c765"}]`,
			},
		},
		{
			name:     "double check calc hash",
			jsonData: DATA,
			hashKey:  "123",
			db:       storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    `[{"id":"PollCount","type":"counter","delta":10,"hash":"5848f030886e828a0481a33dd7386ee4c617c17b4f315645b83dd4b129514639"},{"id":"SomeParam","type":"gauge","value":10.1,"hash":"4b37bfcaa5c419ce160ed8a1247c66a109b086336ac4d7d9df659a86c764a93b"},{"id":"PollCount","type":"counter","delta":20,"hash":"a259b544859dcd20f29fa2d140b6e22e051a200eac13cb37849a91487dcf9821"},{"id":"SomeParam","type":"gauge","value":11.1,"hash":"335bd95edd095efeb1d61bda8bc51f613690b8a6908866a2169a751bcf960933"}]`,
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
