package handler

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sherhan361/monitor/internal/common"
	"github.com/sherhan361/monitor/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

func (h *Handlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.repository.GetAll()
	ggs, _ := json.Marshal(gauges)
	cnts, _ := json.Marshal(counters)
	str := fmt.Sprintf("Gauges: %s\n Counters: %s\n", string(ggs), string(cnts))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.repository.Ping()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	status := checkParams(typ, name, value)
	w.WriteHeader(status)
	if status != http.StatusOK {
		return
	}
	err := h.repository.Set(ctx, typ, name, value)
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

func (h *Handlers) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var input models.Metric
	decodeData := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decodeData.Decode(&input)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := h.repository.GetMetricsByID(input.ID, input.MType, h.cfg.Key)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	js, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
	}

}

func (h *Handlers) CreateMetricsFromJSON(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
		defer r.Body.Close()
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metric := models.Metric{}
	erUnm := json.Unmarshal(body, &metric)
	if erUnm != nil {
		log.Println("err", erUnm)
	}

	if h.cfg.Key != "" {
		log.Println("h.cfg.Key:", h.cfg.Key)
		if !isValidHash(metric, h.cfg.Key) {
			log.Println("err hash:")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err = h.repository.SetMetrics(ctx, &metric)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storMetric, err := h.repository.GetMetricsByID(metric.ID, metric.MType, h.cfg.Key)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if h.cfg.Key != "" {
		storMetric.Hash = common.GetHash(*storMetric, h.cfg.Key)
	}
	js, err := json.Marshal(storMetric)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
	}

}

func (h *Handlers) CreateMetricBatchJSON(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
		defer r.Body.Close()
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var metrics []models.Metric
	erUnm := json.Unmarshal(body, &metrics)
	if erUnm != nil {
		log.Println("err", erUnm)
	}

	if h.cfg.Key != "" {
		log.Println("h.cfg.Key:", h.cfg.Key)
		for _, metric := range metrics {
			if !isValidHash(metric, h.cfg.Key) {
				log.Println("err hash:")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	err = h.repository.SetMetricsBatch(ctx, metrics)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storageMetrics []models.Metric
	for _, metric := range metrics {
		storMetric, err := h.repository.GetMetricsByID(metric.ID, metric.MType, h.cfg.Key)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if h.cfg.Key != "" {
			storMetric.Hash = common.GetHash(*storMetric, h.cfg.Key)
		}
		storageMetrics = append(storageMetrics, *storMetric)
	}

	js, err := json.Marshal(storageMetrics)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
	}

}

func isValidHash(m models.Metric, key string) bool {
	return m.Hash == common.GetHash(m, key)
}
