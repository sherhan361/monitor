package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sherhan361/monitor/internal/server/config"
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"github.com/sherhan361/monitor/internal/server/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("server cfg", cfg)

	strg := GetStor(cfg)

	if cfg.Restore {
		err := strg.RestoreMetrics(cfg.StoreFile)
		if err != nil {
			log.Println(err)
		}
	}

	backuper, err := service.NewBackuper(cfg.StoreInterval, cfg.StoreFile, strg)
	if err != nil {
		log.Fatalln(err)
	}
	go backuper.Run()

	h := handler.NewHandlers(strg, cfg)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, h.Routes()))
}

func GetStor(cfg config.Config) repository.Getter {
	if cfg.DSN == "" {
		strg, err := repository.NewMemoryStorage(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		return strg

	} else {
		strg, err := repository.NewDBStorage(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		return strg
	}
}
