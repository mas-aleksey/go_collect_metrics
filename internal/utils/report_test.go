package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONReport(t *testing.T) {
	statistic := NewStatistic()
	report := NewJSONReport(statistic, "123")

	assert.GreaterOrEqual(t, len(report.Metrics), 32)
	assert.GreaterOrEqual(t, cap(report.Metrics), 32)
}
