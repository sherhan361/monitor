package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sherhan361/monitor/internal/common"
	"github.com/sherhan361/monitor/internal/models"
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

func NewMonitor(pollInterval time.Duration, reportInterval time.Duration, baseURL string, key string) {
	m := &Metrics{
		mutex:    sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	startMonitor(m, pollInterval, reportInterval, baseURL, key)
}

func startMonitor(m *Metrics, pollInterval time.Duration, reportInterval time.Duration, baseURL string, key string) {
	var rtm runtime.MemStats
	var lastSend time.Time
	for {
		<-time.After(pollInterval)
		runtime.ReadMemStats(&rtm)
		updateMetrics(m, &rtm)
		if time.Since(lastSend) >= reportInterval {
			sendReport(m, baseURL, key)
			sendBatchReport(m, baseURL, key)
			lastSend = time.Now()
		}
	}
}

func sendReport(m *Metrics, baseURL string, signKey string) {
	var client = http.Client{}
	url := fmt.Sprintf("http://%s/%s", baseURL, "update")
	contentType := "application/json"

	for key, value := range m.Gauges {
		oneMetric := models.Metric{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}
		if signKey != "" {
			oneMetric.Hash = common.GetHash(oneMetric, signKey)
		}
		fmt.Println("gauge oneMetric:", oneMetric)
		metricJSON, err := json.Marshal(oneMetric)
		if err != nil {
			fmt.Printf("json Gauges Error: %s\n", err)
		}
		resp, err := client.Post(url, contentType, bytes.NewBuffer(metricJSON))
		if err != nil {
			fmt.Printf("Send Gauges Error: %s\n", err)
		} else {
			resp.Body.Close()
		}

	}
	fmt.Println("app.metrics.Counters:", m.Counters)
	for key, value := range m.Counters {
		oneMetric := models.Metric{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
		if signKey != "" {
			oneMetric.Hash = common.GetHash(oneMetric, signKey)
		}
		fmt.Println("counter oneMetric:", oneMetric)
		metricJSON, err := json.Marshal(oneMetric)
		if err != nil {
			fmt.Printf("json Counters Error: %s\n", err)
		}
		resp, err := client.Post(url, contentType, bytes.NewBuffer(metricJSON))
		if err != nil {
			fmt.Printf("Send Counters Error: %s\n", err)
		} else {
			m.Counters["PollCount"] = 0
			resp.Body.Close()
		}
	}
}

func sendBatchReport(m *Metrics, baseURL string, signKey string) {
	var metrics []models.Metric
	var client = http.Client{}
	url := fmt.Sprintf("http://%s/%s", baseURL, "updates")
	contentType := "application/json"

	for key, value := range m.Gauges {
		oneMetric := models.Metric{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}
		if signKey != "" {
			oneMetric.Hash = common.GetHash(oneMetric, signKey)
		}
		metrics = append(metrics, oneMetric)
	}
	for key, value := range m.Counters {
		oneMetric := models.Metric{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
		if signKey != "" {
			oneMetric.Hash = common.GetHash(oneMetric, signKey)
		}
		metrics = append(metrics, oneMetric)
	}
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		fmt.Printf("json Error: %s\n", err)
	}
	fmt.Println("metrics:", metrics)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(metricsJSON))
	if err != nil {
		fmt.Printf("Send Error: %s\n", err)
	} else {
		m.Counters["PollCount"] = 0
		resp.Body.Close()
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
