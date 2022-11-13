package utils

import (
	"math/rand"
	"runtime"
)

type Statistic struct {
	Counter  int64
	RndValue float64
	Rtm      runtime.MemStats
}

func (s *Statistic) Collect() {
	s.Counter++
	s.RndValue = rand.Float64()
	runtime.ReadMemStats(&s.Rtm)
}

func NewStatistic() *Statistic {
	s := &Statistic{
		Counter:  0,
		RndValue: rand.Float64(),
	}
	runtime.ReadMemStats(&s.Rtm)
	return s
}
