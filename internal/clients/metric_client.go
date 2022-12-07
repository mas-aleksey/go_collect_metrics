package clients

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"io"
	"net/http"
	"reflect"
)

type MetricClient struct {
	baseURL string
	client  *http.Client
}

func NewMetricClient(baseURL string) MetricClient {
	return MetricClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (mc MetricClient) SendMetrics(statistic *utils.Statistic) {
	statistic.Mutex.RLock()
	defer statistic.Mutex.RUnlock()
	_, _ = mc.postMetric(utils.NewMetric("counter", "PollCount", utils.ToStr(statistic.Counter)))
	for _, metricName := range utils.RuntimeMetricNames {
		val := reflect.ValueOf(&statistic.Rtm).Elem().FieldByName(metricName).Interface()
		_, _ = mc.postMetric(utils.NewMetric("gauge", metricName, utils.ToStr(val)))
	}
	_, _ = mc.postMetric(utils.NewMetric("gauge", "RandomValue", utils.ToStr(statistic.RndValue)))
	fmt.Println("Report statistic", statistic.Counter)
}

func (mc MetricClient) postMetric(metric utils.Metric) (string, error) {
	endpoint := mc.baseURL + "/update/" + string(metric.Type) + "/" + metric.Name + "/" + metric.Value
	response, err := mc.client.Post(endpoint, "text/plain", nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to complete Post request")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response data")
	}
	return string(body), nil
}
