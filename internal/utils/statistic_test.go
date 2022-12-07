package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStatistic(t *testing.T) {
	statistic := NewStatistic()

	assert.Equal(t, int64(0), statistic.Counter)
	assert.NotNil(t, statistic.RndValue)
	assert.NotNil(t, statistic.Rtm)
}

func TestStatistic_Collect(t *testing.T) {
	statistic := NewStatistic()

	assert.Equal(t, int64(0), statistic.Counter)
	rnd1 := statistic.RndValue
	statistic.Collect()

	rnd2 := statistic.RndValue
	assert.Equal(t, int64(1), statistic.Counter)
	assert.NotEqual(t, rnd1, rnd2)
}
