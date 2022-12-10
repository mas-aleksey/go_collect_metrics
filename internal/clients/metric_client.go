package clients

import (
	"github.com/pkg/errors"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"net/http"
	"sync"
	"time"
)

type MetricClient struct {
	baseURL string
	client  *http.Client
}

func NewMetricClient(baseURL string, timeout time.Duration) MetricClient {
	return MetricClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (mc MetricClient) SendReport(report *utils.Report) error {
	var wg sync.WaitGroup
	errChains := make(chan error, utils.ReportCount)

	for _, metric := range report.Metrics {
		wg.Add(1)
		metric := metric

		go func() {
			defer wg.Done()
			err := mc.postMetric(metric)
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

func (mc MetricClient) postMetric(metric utils.Metric) error {
	endpoint := mc.baseURL + "/update/" + string(metric.Type) + "/" + metric.Name + "/" + metric.Value
	response, err := mc.client.Post(endpoint, "text/plain", nil)
	if err != nil {
		return errors.Wrap(err, "unable to complete Post request")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	return nil
}
