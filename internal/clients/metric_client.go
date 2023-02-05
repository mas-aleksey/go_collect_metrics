package clients

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
	"sync"
	"time"
)

type MetricClient struct {
	*BaseClient
}

func NewMetricClient(baseURL string, timeout time.Duration) *MetricClient {
	return &MetricClient{
		BaseClient: NewBaseClient(baseURL, timeout),
	}
}

func (mc MetricClient) SendJSONReport(report *utils.JSONReport) error {
	var wg sync.WaitGroup
	errChains := make(chan error, utils.ReportCount)

	for _, metric := range report.Metrics {
		wg.Add(1)
		metric := metric

		go func() {
			defer wg.Done()
			body, err := json.Marshal(metric)
			if err != nil {
				errChains <- err
			} else {
				err = mc.postBody(body, false)
				if err != nil {
					errChains <- err
				}
			}
		}()
	}
	wg.Wait()
	close(errChains)

	for err := range errChains {
		if err != nil {
			return err
		}
	}
	return nil
}

func (mc MetricClient) SendBatchJSONReport(report *utils.JSONReport) error {
	body, err := json.Marshal(report.Metrics)
	if err != nil {
		return errors.Wrap(err, "unable to make json")
	}
	return mc.postBody(body, true)
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

func (mc MetricClient) postBody(body []byte, compress bool) error {
	request := Request{
		Method:       http.MethodPost,
		URL:          mc.MakeURL("update/"),
		Headers:      mc.getHeaders(compress),
		Body:         body,
		OkStatusCode: http.StatusOK,
	}
	_, err := mc.DoRequest(request)
	if err != nil {
		return errors.Wrap(err, "unable to complete update metric request")
	}
	return nil
}
