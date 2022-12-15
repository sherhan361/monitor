package main

import (
	"github.com/sherhan361/monitor/internal/server/handler"
	"github.com/sherhan361/monitor/internal/server/repository"
	"log"
	"net/http"
)

func main() {
	strg, err := repository.NewGetter()
	if err != nil {
		log.Fatalln(err)
	}
	h := handler.NewHandlers(strg)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", h.Routes()))
}
