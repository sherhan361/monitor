package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sherhan361/monitor/internal/storage"
)

type Storage interface {
	storage.Getter
}

type Handlers struct {
	storage Storage
}

func NewHandlers(storage storage.Getter) *Handlers {
	return &Handlers{storage: storage}
}

func (h *Handlers) Routes() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/update").Subrouter()
	s.HandleFunc("/{type}/{name}/{value}", h.CreateMetricHandler).Methods("POST")
	return s
}
