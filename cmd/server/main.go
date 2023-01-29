package main

import (
	"fmt"
	"github.com/sherhan361/monitor/internal/server/config"
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"github.com/sherhan361/monitor/internal/server/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("cfg", cfg)

	strg, err := repository.NewMemoryStorage(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if cfg.Restore {
		err = strg.RestoreMetrics(cfg.StoreFile)
		if err != nil {
			log.Println(err)
		}
	}

	backuper, err := service.NewBackuper(cfg.StoreInterval, cfg.StoreFile, strg)
	if err != nil {
		log.Fatalln(err)
	}
	go backuper.Run()

	h := handler.NewHandlers(strg)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, h.Routes()))
}
