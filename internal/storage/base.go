package storage

import (
	"utils"
)

type Repositories interface {
	SaveMetric(metric utils.Metric)
}
