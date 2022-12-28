package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"log"
	"net/http"
)

type Config struct {
	BaseURL string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	strg, err := repository.NewGetter()
	if err != nil {
		log.Fatalln(err)
	}
	h := handler.NewHandlers(strg)
	fmt.Println("cfg.BaseURL:", cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, h.Routes()))
}
