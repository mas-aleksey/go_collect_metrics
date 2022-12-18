package storage

import (
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

type Repositories interface {
	SaveMetric(metric utils.Metric)
}
