package handler

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMonitor(pollInterval int) *Metrics {
	m := &Metrics{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	go startMonitor(m, pollInterval)

	return m
}

func startMonitor(m *Metrics, pollInterval int) {
	var rtm runtime.MemStats
	var interval = time.Duration(pollInterval) * time.Second
	for {
		<-time.After(interval)
		runtime.ReadMemStats(&rtm)
		updateMetrics(m, &rtm)
	}
}

func updateMetrics(m *Metrics, rtm *runtime.MemStats) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Gauges["Alloc"] = float64(rtm.Alloc)
	m.Gauges["BuckHashSys"] = float64(rtm.BuckHashSys)
	m.Gauges["Frees"] = float64(rtm.Frees)
	m.Gauges["GCCPUFraction"] = rtm.GCCPUFraction
	m.Gauges["GCSys"] = float64(rtm.GCSys)
	m.Gauges["HeapAlloc"] = float64(rtm.HeapAlloc)
	m.Gauges["HeapIdle"] = float64(rtm.HeapIdle)
	m.Gauges["HeapInuse"] = float64(rtm.HeapInuse)
	m.Gauges["HeapObjects"] = float64(rtm.HeapObjects)
	m.Gauges["HeapObjects"] = float64(rtm.HeapObjects)
	m.Gauges["HeapReleased"] = float64(rtm.HeapReleased)
	m.Gauges["HeapSys"] = float64(rtm.HeapSys)
	m.Gauges["LastGC"] = float64(rtm.LastGC)
	m.Gauges["Lookups"] = float64(rtm.Lookups)
	m.Gauges["MCacheInuse"] = float64(rtm.MCacheInuse)
	m.Gauges["MCacheSys"] = float64(rtm.MCacheSys)
	m.Gauges["MSpanInuse"] = float64(rtm.MSpanInuse)
	m.Gauges["MSpanSys"] = float64(rtm.MSpanSys)
	m.Gauges["Mallocs"] = float64(rtm.Mallocs)
	m.Gauges["NextGC"] = float64(rtm.NextGC)
	m.Gauges["NumForcedGC"] = float64(rtm.NumForcedGC)
	m.Gauges["NumGC"] = float64(rtm.NumGC)
	m.Gauges["OtherSys"] = float64(rtm.OtherSys)
	m.Gauges["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	m.Gauges["StackInuse"] = float64(rtm.StackInuse)
	m.Gauges["StackSys"] = float64(rtm.StackSys)
	m.Gauges["Sys"] = float64(rtm.Sys)
	m.Gauges["TotalAlloc"] = float64(rtm.TotalAlloc)
	m.Gauges["RandomValue"] = rand.Float64()

	m.Counters["PollCount"] = m.Counters["PollCount"] + 1
}
