package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/sherhan361/monitor/internal/server/repository"
)

type Storage interface {
	repository.Getter
}

type Handlers struct {
	storage Storage
}

func NewHandlers(storage repository.Getter) *Handlers {
	return &Handlers{storage: storage}
}

func (h *Handlers) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.CreateMetricHandler)
	})
	return r
}
