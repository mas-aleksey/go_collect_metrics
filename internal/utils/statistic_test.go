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
	stat1 := *statistic
	statistic.Collect()
	stat2 := *statistic

	assert.Equal(t, int64(0), stat1.Counter)
	assert.Equal(t, int64(1), stat2.Counter)
	assert.NotEqual(t, stat1, stat2)
	assert.NotEqual(t, stat1.Counter, stat2.Counter)
	assert.NotEqual(t, stat1.RndValue, stat2.RndValue)
}
