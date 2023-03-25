package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sherhan361/monitor/internal/server/config"
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"github.com/sherhan361/monitor/internal/server/service"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("server cfg", cfg)

	db, err := sql.Open("pgx", "postgres://postgres:example@localhost:5432/monitor")
	if err != nil {
		fmt.Println("err:", err)
	}
	defer db.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("ping!")

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

	h := handler.NewHandlers(strg, cfg)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, h.Routes()))
}
