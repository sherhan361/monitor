package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sherhan361/monitor/internal/agent/handler"
	"log"
	"os"
	"time"
)

type Config struct {
	User           string        `env:"USER"`
	BaseURL        string        `env:"ADDRESS" envDefault:"http://127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	u := os.Getenv("ADDRESS")
	fmt.Println("u", u)
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cfg:", cfg)
	handler.NewMonitor(cfg.PollInterval, cfg.ReportInterval, cfg.BaseURL)
}
