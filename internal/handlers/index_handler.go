package handlers

import (
	"github.com/tiraill/go_collect_metrics/internal/storage"
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

// MetricData - структура для данных метрики
type MetricData struct {
	Name  string
	Value string
}

// TemplateData - структура для шаблона данных с метриками
type TemplateData struct {
	Metrics []MetricData
}

// IndexHandler - метод для получения HTML страници со списком всех метрик
// GET /.
func IndexHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data := TemplateData{Metrics: make([]MetricData, 0)}
		metrics, err := db.GetAllMetrics(ctx)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		for _, metric := range metrics {
			data.Metrics = append(data.Metrics, MetricData{Name: metric.ID, Value: metric.ValueString()})
		}
		w.Header().Set("content-type", "text/html")
		fpT, _ := template.New("metrics").Parse(pageTemp)
		err = fpT.Execute(w, data)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
