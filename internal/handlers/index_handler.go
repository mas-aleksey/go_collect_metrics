package handlers

import (
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"html/template"
	"net/http"
)

var pageTemp = `<!DOCTYPE html>
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
      	{{ with .Metrics }}
			{{ range . }}
      			<tr>
              		<td>{{ .Name }}</td>
              		<td>{{ .Value }}</td>
      			</tr>
			{{ end }} 
      	{{ end }}
    </table>
  </body>
</html>`

type MetricData struct {
	Name  string
	Value string
}

type TemplateData struct {
	Metrics []MetricData
}

func IndexHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := TemplateData{Metrics: make([]MetricData, 0)}
		for k, v := range storage.GaugeMetrics {
			data.Metrics = append(data.Metrics, MetricData{Name: k, Value: utils.ToStr(v)})
		}
		for k, v := range storage.CounterMetrics {
			data.Metrics = append(data.Metrics, MetricData{Name: k, Value: utils.ToStr(v)})
		}
		fpT, _ := template.New("metrics").Parse(pageTemp)
		err := fpT.Execute(w, data)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}
