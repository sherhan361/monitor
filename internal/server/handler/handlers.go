package handler

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"

	"github.com/sherhan361/monitor/internal/models"
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

	metric, err := h.repository.GetMetricsByID(input.ID, input.MType)
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
		fmt.Println("err", erUnm)
	}

	err = h.repository.SetMetrics(&metric)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storMetric, err := h.repository.GetMetricsByID(metric.ID, metric.MType)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
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
