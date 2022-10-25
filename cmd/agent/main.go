package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
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
		metrics := *m

		v := reflect.ValueOf(metrics)
		types := v.Type()
		values := make([]interface{}, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			values[i] = v.Field(i).Interface()

			baseUrl := "http://127.0.0.1:8080/update/"
			valueType := "gauge"
			if reflect.TypeOf(values[i]).Name() == "int64" {
				valueType = "counter"
			}
			testUrl := fmt.Sprintf("%v%v/%v/%v", baseUrl, valueType, types.Field(i).Name, values[i])
			fmt.Println("testUrl:", testUrl)
			metricsValue, _ := json.Marshal(metrics)

			//_, err := http.Post(testUrl, "text/plain", bytes.NewBuffer(metricsValue))
			//if err != nil {
			//	fmt.Println("err:", err)
			//}

			req, err := http.NewRequest(http.MethodPost, testUrl, bytes.NewBuffer(metricsValue))
			if err != nil {
				fmt.Printf("client: could not create request: %s\n", err)
				os.Exit(1)
			}
			req.Header.Set("Content-Type", "text/plain")
			client := http.Client{
				Timeout: 30 * time.Second,
			}
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("client: error making http request: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("res:", res)

		}
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
	}

}

func main() {
	pollInterval := 2
	reportInterval := 10
	var m Metrics

	go Monitor(&m, pollInterval)
	go ReportSender(&m, reportInterval)

	//var interval = time.Duration(reportInterval) * time.Second
	//for {
	//	<-time.After(interval)
	//	fmt.Println("m:", m)
	//}
}
