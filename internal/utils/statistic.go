package utils

import (
	"fmt"
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
	fmt.Println("Collect statistic", s.Counter)
}

func NewStatistic() *Statistic {
	s := &Statistic{
		Counter:  0,
		RndValue: rand.Float64(),
	}
	runtime.ReadMemStats(&s.Rtm)
	return s
}
