package utils

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"runtime"
	"sync"
)

type Statistic struct {
	Counter  int64
	RndValue float64
	MemStat  *mem.VirtualMemoryStat
	CpuCount int
	Rtm      runtime.MemStats
	Mutex    sync.RWMutex
}

func NewStatistic() *Statistic {
	s := &Statistic{
		Counter:  0,
		RndValue: rand.Float64(),
	}
	s.MemStat, _ = mem.VirtualMemory()
	s.CpuCount, _ = cpu.Counts(false)
	runtime.ReadMemStats(&s.Rtm)
	return s
}

func (s *Statistic) CollectRuntime() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Counter++
	s.RndValue = rand.Float64()
	runtime.ReadMemStats(&s.Rtm)
	log.Println("Collect runtime statistic", s.Counter)
}

func (s *Statistic) CollectMemCpu() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.MemStat, _ = mem.VirtualMemory()
	s.CpuCount, _ = cpu.Counts(false)
	log.Println("Collect mem cpu statistic")
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
		MemStat:  s.MemStat,
		CpuCount: s.CpuCount,
		Rtm:      s.Rtm,
	}
}
