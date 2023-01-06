package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"github.com/sherhan361/monitor/internal/server/service"
	"log"
	"net/http"
	"time"
)

type Config struct {
	BaseURL       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

type ArgConfig struct {
	BaseURL       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	var argCfg ArgConfig
	flag.StringVar(&argCfg.BaseURL, "a", "127.0.0.1:8080", "host:port")
	flag.DurationVar(&argCfg.StoreInterval, "i", time.Duration(300*time.Second), "backup interval")
	flag.StringVar(&argCfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "filename to backup")
	flag.BoolVar(&argCfg.Restore, "r", true, "is restore enabled")
	flag.Parse()
	if cfg.BaseURL == "" {
		cfg.BaseURL = argCfg.BaseURL
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = argCfg.StoreInterval
	}
	if cfg.StoreFile == "" {
		cfg.StoreFile = argCfg.StoreFile
	}
	if !cfg.Restore {
		cfg.Restore = argCfg.Restore
	}

	strg, err := repository.NewMemoryStorage()
	if err != nil {
		log.Fatalln(err)
	}

	if cfg.Restore {
		err = strg.RestoreMetrics(cfg.StoreFile)
		if err != nil {
			log.Println(err)
		}
	}

	producer, err := service.NewBackuper(cfg.StoreInterval, cfg.StoreFile, strg)
	if err != nil {
		log.Fatalln(err)
	}
	go producer.Run()

	h := handler.NewHandlers(strg)
	fmt.Println("cfg", cfg)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, h.Routes()))
}
