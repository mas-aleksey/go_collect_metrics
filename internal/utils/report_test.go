package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONReport(t *testing.T) {
	statistic := NewStatistic()
	report := NewJSONReport(statistic, "123")

	assert.GreaterOrEqual(t, len(report.Metrics), ReportCount)
	assert.GreaterOrEqual(t, cap(report.Metrics), ReportCount)
}
