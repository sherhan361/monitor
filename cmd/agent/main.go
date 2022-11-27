package main

import (
	"github.com/sherhan361/monitor/internal/agent/handler"
)

func main() {
	var pollInterval = 2
	var reportInterval = 10
	var baseURL = "http://127.0.0.1:8080"

	handler.NewMonitor(pollInterval, reportInterval, baseURL)
}
