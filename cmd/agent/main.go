package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type Metrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapObjects   float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64

	PollCount   int64
	RandomValue float64
}

func ReportSender(m *Metrics, reportInterval int) {
	var interval = time.Duration(reportInterval) * time.Second
	for {
		<-time.After(interval)
		fmt.Println("m:", m)
		url := "http://127.0.0.1:8080"
		metricsValue, _ := json.Marshal(m)
		getUrl(m)
		_, err := http.Post(url, "application/json", bytes.NewBuffer(metricsValue))
		if err != nil {
			fmt.Println("err:", err)
		}
	}
}

func getUrl(m *Metrics) {
	t := *m
	v := reflect.ValueOf(t)
	fmt.Println("RES:", v)
	typeOfS := v.Type()
	fmt.Println("typeOfS:", typeOfS)
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v Type: %s\n", typeOfS.Field(i).Name, v.Field(i).Interface(), v.Field(i).Type())
	}

}

func Monitor(m *Metrics, pollInterval int) {
	var rtm runtime.MemStats
	var interval = time.Duration(pollInterval) * time.Second
	var PollCount int64 = 0
	for {
		<-time.After(interval)
		runtime.ReadMemStats(&rtm)

		m.Alloc = float64(rtm.Alloc)
		m.BuckHashSys = float64(rtm.BuckHashSys)
		m.Frees = float64(rtm.Frees)
		m.GCCPUFraction = rtm.GCCPUFraction
		m.GCSys = float64(rtm.GCSys)
		m.HeapAlloc = float64(rtm.HeapAlloc)
		m.HeapIdle = float64(rtm.HeapIdle)
		m.HeapObjects = float64(rtm.HeapObjects)
		m.HeapSys = float64(rtm.HeapSys)
		m.LastGC = float64(rtm.LastGC)
		m.Lookups = float64(rtm.Lookups)
		m.MCacheInuse = float64(rtm.MCacheInuse)
		m.MCacheSys = float64(rtm.MCacheSys)
		m.MSpanInuse = float64(rtm.MSpanInuse)
		m.MSpanSys = float64(rtm.MSpanSys)
		m.Mallocs = float64(rtm.Mallocs)
		m.NextGC = float64(rtm.NextGC)
		m.NumForcedGC = float64(rtm.NumForcedGC)
		m.NumGC = float64(rtm.NumGC)
		m.OtherSys = float64(rtm.OtherSys)
		m.PauseTotalNs = float64(rtm.PauseTotalNs)
		m.StackInuse = float64(rtm.StackInuse)
		m.StackSys = float64(rtm.StackSys)
		m.Sys = float64(rtm.Sys)
		m.TotalAlloc = float64(rtm.TotalAlloc)

		m.PollCount = PollCount
		PollCount++
		m.RandomValue = rand.Float64()
		fmt.Println("Metrics:", m)
	}

}

func main() {
	pollInterval := 1
	reportInterval := 2
	var m Metrics

	go Monitor(&m, pollInterval)
	go ReportSender(&m, reportInterval)

	var interval = time.Duration(reportInterval) * time.Second
	for {
		<-time.After(interval)
		fmt.Println("m:", m)
	}
}
