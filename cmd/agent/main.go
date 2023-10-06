package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sherhan361/monitor/internal/agent/handler"
	"log"
	"time"
)

type Config struct {
	BaseURL        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

type ArgConfig struct {
	BaseURL        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
	RateLimit      int
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	var argCfg ArgConfig
	flag.StringVar(&argCfg.BaseURL, "a", "127.0.0.1:8080", "host:port")
	flag.DurationVar(&argCfg.ReportInterval, "r", time.Duration(10*time.Second), "report interval")
	flag.DurationVar(&argCfg.PollInterval, "p", time.Duration(2*time.Second), "poll interval")
	flag.StringVar(&argCfg.Key, "k", "", "sign key")
	flag.IntVar(&argCfg.RateLimit, "l", 1, "RATE LIMIT")
	flag.Parse()

	log.Println("agent argCfg.Key:", argCfg.Key)
	if cfg.BaseURL == "" {
		cfg.BaseURL = argCfg.BaseURL
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = argCfg.ReportInterval
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = argCfg.PollInterval
	}
	log.Println("agent cfg.Key:", cfg.Key)
	if cfg.Key == "" {
		cfg.Key = argCfg.Key
	}
	if cfg.RateLimit == 0 {
		cfg.RateLimit = argCfg.RateLimit
	}

	log.Println("agent cfg:", cfg)
	handler.NewMonitor(cfg.PollInterval, cfg.ReportInterval, cfg.BaseURL, cfg.Key, cfg.RateLimit)
}
