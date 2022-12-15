package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (h *Handlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.repository.GetAll()
	ggs, _ := json.Marshal(gauges)
	cnts, _ := json.Marshal(counters)
	str := fmt.Sprintf("Gauges: %s\n Counters: %s\n", string(ggs), string(cnts))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handlers) GetMetric(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	if typ == "" && name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	str, err := h.repository.Get(typ, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(str))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handlers) CreateMetric(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	status := checkParams(typ, name, value)
	w.WriteHeader(status)
	if status != http.StatusOK {
		return
	}
	err := h.repository.Set(typ, name, value)
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
	if name == "none" || value == "none" {
		return http.StatusBadRequest
	}
	if typ != "gauge" && typ != "counter" {
		return http.StatusNotImplemented
	}
	return http.StatusOK
}
