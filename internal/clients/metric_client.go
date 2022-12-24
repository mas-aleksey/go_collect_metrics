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
			err := mc.postJSONMetric(metric)
			errChains <- err
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

func (mc MetricClient) getHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return headers
}

func (mc MetricClient) postJSONMetric(metric utils.JSONMetric) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return errors.Wrap(err, "unable to make json")
	}
	request := Request{
		Method:       http.MethodPost,
		URL:          mc.MakeURL("update/"),
		Headers:      mc.getHeaders(),
		Body:         body,
		OkStatusCode: http.StatusOK,
	}
	_, err = mc.DoRequest(request)
	if err != nil {
		return errors.Wrap(err, "unable to complete update metric request")
	}
	return nil
}
