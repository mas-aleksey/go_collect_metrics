package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	statistic.CollectRuntime()
	statistic.CollectMemCPU()

	rnd2 := statistic.RndValue
	assert.Equal(t, int64(1), statistic.Counter)
	assert.NotEqual(t, rnd1, rnd2)
}

func TestStatistic_ResetCounter(t *testing.T) {
	statistic := NewStatistic()
	statistic.CollectRuntime()
	statistic.CollectRuntime()
	statistic.CollectMemCPU()
	assert.Equal(t, int64(2), statistic.Counter)
	statistic.ResetCounter()
	assert.Equal(t, int64(0), statistic.Counter)
}

func TestStatistic_Copy(t *testing.T) {
	statistic := NewStatistic()
	statistic.CollectRuntime()
	statistic.CollectMemCPU()
	statCopy := statistic.Copy()

	assert.Equal(t, statCopy.Counter, statistic.Counter)
	assert.Equal(t, statCopy.RndValue, statistic.RndValue)
	assert.Equal(t, statCopy.Rtm, statistic.Rtm)
	assert.Equal(t, statCopy, statistic)
	assert.NotSame(t, statCopy, statistic)
}
