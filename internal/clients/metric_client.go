package clients

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

// MetricClient - струкрура описывает клиента для отпарвки метрик
type MetricClient struct {
	*BaseClient
}

// NewMetricClient - метод для создания клиента отпарвки метрик
func NewMetricClient(baseURL string, timeout time.Duration, rateLimit int) *MetricClient {
	return &MetricClient{
		BaseClient: NewBaseClient(baseURL, timeout, rateLimit),
	}
}

// SendJSONReport - метод для отправки отчета в формате JSON
// для отправки каждой метрики выполняется отдельный API запрос
//
// Deprecated: используйте метод SendBatchJSONReport
func (mc MetricClient) SendJSONReport(report *utils.JSONReport) error {
	ctx, cancel := context.WithCancel(context.Background())
	g, _ := errgroup.WithContext(ctx)

	metricCh := make(chan utils.JSONMetric, 33)

	for i := 0; i < mc.rateLimit; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case metric := <-metricCh:
					body, err := json.Marshal(metric)
					if err != nil {
						return err
					} else {
						err = mc.postBody(body, "update/", false)
						if err != nil {
							return err
						}
					}
				}
			}
		})
	}
	g.Go(func() error {
		for _, metric := range report.Metrics {
			//m := &metric
			metricCh <- metric
		}
		close(metricCh)
		cancel()
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// SendBatchJSONReport - метод для отправки отчета в формате JSON
// для отправки всех метрик будет выполнен один API запрос
func (mc MetricClient) SendBatchJSONReport(report *utils.JSONReport) error {
	body, err := json.Marshal(report.Metrics)
	if err != nil {
		return errors.Wrap(err, "unable to make json")
	}
	return mc.postBody(body, "updates/", true)
}

func (mc MetricClient) getHeaders(compress bool) map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept-Encoding"] = "gzip"
	if compress {
		headers["Content-Encoding"] = "gzip"
	}
	return headers
}

func (mc MetricClient) postBody(body []byte, url string, compress bool) error {
	request := Request{
		Method:       http.MethodPost,
		URL:          mc.MakeURL(url),
		Headers:      mc.getHeaders(compress),
		Body:         body,
		OkStatusCode: http.StatusOK,
	}
	_, err := mc.DoRequest(&request)
	if err != nil {
		return errors.Wrap(err, "unable to complete update metric request")
	}
	return nil
}
