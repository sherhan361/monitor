package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sherhan361/monitor/internal/server/repository"
)

type Repository interface {
	repository.Getter
}

type Handlers struct {
	repository Repository
}

func NewHandlers(repository repository.Getter) *Handlers {
	return &Handlers{repository: repository}
}

func (h *Handlers) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/", h.GetAllMetrics)
	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", h.GetMetric)
		r.Post("/", h.GetMetricsJSON)
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.CreateMetric)
		r.Post("/", h.CreateMetricsFromJSON)
	})
	return r
}
