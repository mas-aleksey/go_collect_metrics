package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONReport(t *testing.T) {
	statistic := NewStatistic()
	report := NewJSONReport(statistic)

	assert.Equal(t, len(report.Metrics), ReportCount)
	assert.Equal(t, cap(report.Metrics), ReportCount)
}
