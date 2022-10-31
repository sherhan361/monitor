package main

import (
	"github.com/sherhan361/monitor/internal/handlers"
	"github.com/sherhan361/monitor/internal/storage"
	"log"
	"net/http"
)

func main() {
	strg, err := storage.NewGetter()
	if err != nil {
		log.Fatalln(err)
	}
	h := handlers.NewHandlers(strg)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", h.Routes()))
}
