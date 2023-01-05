package main

import (
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
	BaseURL       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
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
