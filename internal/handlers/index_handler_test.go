package handlers

import (
	"compress/gzip"
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

var emptyPage = `<!DOCTYPE html>
<html>
  <head>
    <title>Metrics</title>
  </head>
  <body>
    <p>
      List of metrics.
    </p>
    <table>
      	<tr>
          <td>Name</td>
          <td>Value</td>
    	</tr>
      	
    </table>
  </body>
</html>`

var fillPage = `<!DOCTYPE html>
<html>
  <head>
    <title>Metrics</title>
  </head>
  <body>
    <p>
      List of metrics.
    </p>
    <table>
      	<tr>
          <td>Name</td>
          <td>Value</td>
    	</tr>
      	
			
      			<tr>
              		<td>Alloc</td>
              		<td>111.222</td>
      			</tr>
			
      			<tr>
              		<td>PollCount</td>
              		<td>333</td>
      			</tr>
			 
      	
    </table>
  </body>
</html>`

func TestGetIndexMetricHandler(t *testing.T) {
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
		name string
		db   storage.Storage
		want want
	}{
		{
			name: "check 200 empty metrics",
			db:   storage.NewStorage(&utils.StorageConfig{}),
			want: want{
				statusCode: 200,
				message:    emptyPage,
			},
		},
		{
			name: "check 200 fill some metrics",
			db:   testStorage,
			want: want{
				statusCode: 200,
				message:    fillPage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.db, utils.ServerConfig{Address: "adr", HashKey: "key"})
			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			result, err := client.Get(ts.URL)
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

func TestGetCompressedPage(t *testing.T) {

	testStorage := storage.NewStorage(&utils.StorageConfig{})
	_, _ = testStorage.UpdateJSONMetrics(context.Background(), []utils.JSONMetric{
		utils.NewGaugeJSONMetric("Alloc", 111.222),
		utils.NewCounterJSONMetric("PollCount", 333),
	})

	r := GetRouter(testStorage, utils.ServerConfig{})
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	result, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "gzip", result.Header.Get("Content-Encoding"))
	assert.Equal(t, "225", result.Header.Get("Content-Length"))

	gzipReader, err := gzip.NewReader(result.Body)
	require.NoError(t, err)

	resBody, err := io.ReadAll(gzipReader)
	require.NoError(t, err)

	err = result.Body.Close()
	require.NoError(t, err)
	err = gzipReader.Close()
	require.NoError(t, err)

	assert.Equal(t, fillPage, string(resBody))
}
