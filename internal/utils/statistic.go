package utils

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"runtime"
	"sync"
)

// Statistic - струкрура для хранения метрик.
type Statistic struct {
	Counter        int64
	RndValue       float64
	MemStat        *mem.VirtualMemoryStat
	CPUUtilization []float64
	Rtm            *runtime.MemStats
	Mutex          sync.RWMutex
}

func getRtm() *runtime.MemStats {
	var buf runtime.MemStats
	runtime.ReadMemStats(&buf)
	return &buf
}

func getMemStat() *mem.VirtualMemoryStat {
	memStat, _ := mem.VirtualMemory()
	return memStat
}

func getCpuStat() []float64 {
	CPUUtilization, _ := cpu.Percent(0, true)
	return CPUUtilization
}

// NewStatistic - метод для создания струкруры Statistic.
func NewStatistic() *Statistic {
	s := &Statistic{
		Counter:        0,
		RndValue:       rand.Float64(),
		MemStat:        getMemStat(),
		CPUUtilization: getCpuStat(),
		Rtm:            getRtm(),
	}
	return s
}

// CollectRuntime - метод для сбора Runtime метрик.
func (s *Statistic) CollectRuntime() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Counter++
	s.RndValue = rand.Float64()
	s.Rtm = getRtm()
	log.Println("Collect runtime statistic", s.Counter)
}

// CollectMemCPU - метод для сбора Memory и CPU метрик.
func (s *Statistic) CollectMemCPU() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.MemStat = getMemStat()
	s.CPUUtilization = getCpuStat()
	log.Println("Collect mem cpu statistic")
}

// ResetCounter - метод для сброса счетчика Counter.
func (s *Statistic) ResetCounter() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Counter = 0
}

// Copy - метод для создания копии структуры Statistic.
func (s *Statistic) Copy() *Statistic {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return &Statistic{
		Counter:        s.Counter,
		RndValue:       s.RndValue,
		MemStat:        s.MemStat,
		CPUUtilization: s.CPUUtilization,
		Rtm:            s.Rtm,
	}
}
