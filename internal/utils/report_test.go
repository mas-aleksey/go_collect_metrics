package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReport(t *testing.T) {
	statistic := NewStatistic()
	report := NewReport(statistic)

	assert.Equal(t, len(report.Metrics), ReportCount)
	assert.Equal(t, cap(report.Metrics), ReportCount)
}
