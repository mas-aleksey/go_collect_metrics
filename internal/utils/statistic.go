package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

type Statistic struct {
	Counter  int64
	RndValue float64
	Rtm      runtime.MemStats
	Mutex    sync.RWMutex
}

func NewStatistic() *Statistic {
	s := &Statistic{
		Counter:  0,
		RndValue: rand.Float64(),
	}
	runtime.ReadMemStats(&s.Rtm)
	return s
}

func (s *Statistic) Collect() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Counter++
	s.RndValue = rand.Float64()
	runtime.ReadMemStats(&s.Rtm)
	fmt.Println("Collect statistic", s.Counter)
}

func (s *Statistic) ResetCounter() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Counter = 0
}

func (s *Statistic) Copy() *Statistic {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return &Statistic{
		Counter:  s.Counter,
		RndValue: s.RndValue,
		Rtm:      s.Rtm,
	}
}
