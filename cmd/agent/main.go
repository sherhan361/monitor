package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}

func main() {
	pollInterval := 2
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	fmt.Println("Alloc", rtm.Alloc)
	fmt.Println("BuckHashSys", rtm.BuckHashSys)
	fmt.Println("rtm.Frees:", rtm.Frees)
	fmt.Println("rtm.GCCPUFraction:", rtm.GCCPUFraction)
	fmt.Println("GCSys:", rtm.GCSys)
	fmt.Println("HeapAlloc:", rtm.HeapAlloc)
	fmt.Println("HeapIdle:", rtm.HeapIdle)
	fmt.Println("HeapObjects:", rtm.HeapObjects)
	fmt.Println("HeapSys:", rtm.HeapSys)
	fmt.Println("LastGC:", rtm.LastGC)
	fmt.Println("Lookups:", rtm.Lookups)
	fmt.Println("MCacheInuse:", rtm.MCacheInuse)
	fmt.Println("MCacheSys:", rtm.MCacheSys)
	fmt.Println("MSpanInuse:", rtm.MSpanInuse)
	fmt.Println("MSpanSys:", rtm.MSpanSys)
	fmt.Println("Mallocs:", rtm.Mallocs)
	fmt.Println("NextGC:", rtm.NextGC)
	fmt.Println("NumForcedGC:", rtm.NumForcedGC)
	fmt.Println("NumGC:", rtm.NumGC)
	fmt.Println("OtherSys:", rtm.OtherSys)
	fmt.Println("PauseTotalNs:", rtm.PauseTotalNs)
	fmt.Println("StackInuse:", rtm.StackInuse)
	fmt.Println("StackSys:", rtm.StackSys)
	fmt.Println("Sys:", rtm.Sys)
	fmt.Println("TotalAlloc:", rtm.TotalAlloc)

	NewMonitor(pollInterval)
}
