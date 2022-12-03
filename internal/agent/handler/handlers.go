package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMonitor(pollInterval int, reportInterval int, baseURL string) {
	m := &Metrics{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	startMonitor(m, time.Duration(pollInterval)*time.Second, reportInterval, baseURL)
}

func startMonitor(m *Metrics, pollInterval time.Duration, reportInterval int, baseURL string) {
	var rtm runtime.MemStats
	var lastSend time.Time
	for {
		<-time.After(pollInterval)
		runtime.ReadMemStats(&rtm)
		updateMetrics(m, &rtm)
		dif := int(time.Since(lastSend) / time.Second)
		if dif >= reportInterval {
			sendReport(m, baseURL)
			lastSend = time.Now()
		}
	}
}

func sendReport(m *Metrics, baseURL string) {
	var client = http.Client{}

	for key, value := range m.Gauges {
		url := fmt.Sprintf("%s/%s/%s/%s/%.3f", baseURL, "update", "gauge", key, value)
		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Printf("Send Gauges Error: %s\n", err)
		} else {
			resp.Body.Close()
		}

	}
	fmt.Println("app.metrics.Counters:", m.Counters)
	for key, value := range m.Counters {
		url := fmt.Sprintf("%s/%s/%s/%s/%d", baseURL, "update", "counter", key, value)
		resp, err := client.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Printf("Send Counters Error: %s\n", err)
		} else {
			m.Counters["PollCount"] = 0
			resp.Body.Close()
		}
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
