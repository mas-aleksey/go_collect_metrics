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
	testMemStorage := storage.NewMemStorage()
	testMemStorage.GaugeMetrics["Alloc"] = 111.222
	testMemStorage.CounterMetrics["PollCount"] = 333

	tests := []struct {
		name       string
		memStorage *storage.MemStorage
		want       want
	}{
		{
			name:       "check 200 empty metrics",
			memStorage: storage.NewMemStorage(),
			want: want{
				statusCode: 200,
				message:    emptyPage,
			},
		},
		{
			name:       "check 200 fill some metrics",
			memStorage: testMemStorage,
			want: want{
				statusCode: 200,
				message:    fillPage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := GetRouter(tt.memStorage)
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
