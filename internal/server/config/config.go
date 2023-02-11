package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
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

func GetConfig() Config {
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
	return cfg
}
