package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sherhan361/monitor/internal/common"
	"github.com/sherhan361/monitor/internal/models"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	*sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMonitor(pollInterval time.Duration, reportInterval time.Duration, baseURL string, key string, rateLimit int) {
	m := &Metrics{
		RWMutex:  &sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	startMonitor(m, pollInterval, reportInterval, baseURL, key, rateLimit)
}

type Job struct {
	URL         string
	contentType string
	metricsJSON []byte
}

func makePostRequest(url string, contentType string, metricsJSON []byte) {
	var client = http.Client{}
	log.Println("making request to", url)

	resp, err := client.Post(url, contentType, bytes.NewBuffer(metricsJSON))
	if err != nil {
		log.Printf("Send Gauges Error: %s\n", err)
	} else {
		resp.Body.Close()
	}

	log.Println("success", url)
}

func startMonitor(m *Metrics, pollInterval time.Duration, reportInterval time.Duration, baseURL string, key string, rateLimit int) {
	var lastSend time.Time
	wgRefresh := sync.WaitGroup{}

	for {
		<-time.After(pollInterval)
		wgRefresh.Add(2)

		go func() {
			defer wgRefresh.Done()
			updateMetrics(m)
		}()

		go func() {
			updateExtraMetrics(m)
			defer wgRefresh.Done()
		}()
		if time.Since(lastSend) >= reportInterval {
			wgRefresh.Wait()
			jobCh := make(chan *Job)
			for i := 0; i < rateLimit; i++ {
				go func() {
					for job := range jobCh {
						makePostRequest(job.URL, job.contentType, job.metricsJSON)
					}
				}()
			}

			go func() {
				sendReport(m, baseURL, key, jobCh)
				sendBatchReport(m, baseURL, key, jobCh)
			}()
			lastSend = time.Now()
		}
	}
}

func sendReport(m *Metrics, baseURL string, signKey string, jobCh chan *Job) {
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
		log.Println("gauge oneMetric:", oneMetric)
		metricJSON, err := json.Marshal(oneMetric)
		if err != nil {
			log.Printf("json Gauges Error: %s\n", err)
		}

		job := &Job{URL: url, contentType: contentType, metricsJSON: metricJSON}
		jobCh <- job

	}
	log.Println("app.metrics.Counters:", m.Counters)
	for key, value := range m.Counters {
		oneMetric := models.Metric{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
		if signKey != "" {
			oneMetric.Hash = common.GetHash(oneMetric, signKey)
		}
		log.Println("counter oneMetric:", oneMetric)
		metricJSON, err := json.Marshal(oneMetric)
		if err != nil {
			log.Printf("json Counters Error: %s\n", err)
		}
		job := &Job{URL: url, contentType: contentType, metricsJSON: metricJSON}
		jobCh <- job

	}
}

func sendBatchReport(m *Metrics, baseURL string, signKey string, jobCh chan *Job) {
	var metrics []models.Metric
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
		log.Printf("json Error: %s\n", err)
	}
	log.Println("metrics:", metrics)

	job := &Job{URL: url, contentType: contentType, metricsJSON: metricsJSON}
	jobCh <- job

}

func updateMetrics(m *Metrics) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Lock()
	defer m.Unlock()

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

func updateExtraMetrics(m *Metrics) {
	metrics, _ := mem.VirtualMemory()
	m.Lock()
	defer m.Unlock()

	m.Gauges["TotalMemory"] = float64(metrics.Total)
	m.Gauges["FreeMemory"] = float64(metrics.Free)

	percentageCPU, _ := cpu.Percent(0, true)

	for i, currentPercentageCPU := range percentageCPU {
		metricName := fmt.Sprintf("CPUutilization%v", i)
		m.Gauges[metricName] = float64(currentPercentageCPU)
	}
}
