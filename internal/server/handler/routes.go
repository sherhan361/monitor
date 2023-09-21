package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sherhan361/monitor/internal/server/config"
	"github.com/sherhan361/monitor/internal/server/repository"
)

type Repository interface {
	repository.Getter
}

type Handlers struct {
	repository Repository
	cfg        config.Config
}

func NewHandlers(repository repository.Getter, config config.Config) *Handlers {
	return &Handlers{
		repository: repository,
		cfg:        config,
	}
}

func (h *Handlers) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(GzipCompress)

	r.Get("/", h.GetAllMetrics)
	r.Get("/ping", h.Ping)
	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", h.GetMetric)
		r.Post("/", h.GetMetricsJSON)
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.CreateMetric)
		r.Post("/", h.CreateMetricsFromJSON)
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.CreateMetricBatchJSON)
	})
	return r
}
