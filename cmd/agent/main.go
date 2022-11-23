package main

import (
	"flag"
	"fmt"
	"github.com/sherhan361/monitor/internal/agent/handler"
	"net/http"
	"time"
)

type configuration struct {
	pollInterval   int
	reportInterval int
	baseUrl        string
}

type application struct {
	config  configuration
	metrics *handler.Metrics
}

func (app *application) ReportSender() {
	cfg := app.config
	var client = http.Client{}
	var interval = time.Duration(cfg.reportInterval) * time.Second
	for {
		<-time.After(interval)

		for key, value := range app.metrics.Gauges {
			url := fmt.Sprintf("%s/%s/%s/%s/%.3f", cfg.baseUrl, "update", "gauge", key, value)
			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Printf("Gauges Error: %s", err)
			}
			resp.Body.Close()
		}
		for key, value := range app.metrics.Counters {
			url := fmt.Sprintf("%s/%s/%s/%s/%d", cfg.baseUrl, "update", "counter", key, value)
			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				fmt.Printf("Counters Error: %s", err)
			}
			resp.Body.Close()
		}

	}
}

func main() {
	var cfg configuration
	flag.IntVar(&cfg.pollInterval, "pollInterval", 2, "poll interval")
	flag.IntVar(&cfg.reportInterval, "reportInterval", 10, "report interval")
	flag.StringVar(&cfg.baseUrl, "baseUrl", "http://127.0.0.1:8080", "url")
	flag.Parse()

	app := &application{
		config:  cfg,
		metrics: handler.NewMonitor(cfg.pollInterval),
	}
	app.ReportSender()
}
