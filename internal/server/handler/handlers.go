package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	status := checkParams(typ, name, value)
	w.WriteHeader(status)
	if status != http.StatusOK {
		return
	}

	err := h.storage.Set(typ, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func checkParams(typ string, name string, value string) int {
	if typ == "" || name == "" || value == "" {
		return http.StatusBadRequest
	}
	if typ == "none" || name == "none" || value == "none" {
		return http.StatusBadRequest
	}
	if typ != "gauge" && typ != "counter" {
		return http.StatusNotImplemented
	}
	return http.StatusOK
}
