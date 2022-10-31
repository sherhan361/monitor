package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handlers) CreateMetricHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("vars:", vars["name"])
	typ := vars["type"]
	name := vars["name"]
	value := vars["value"]

	if typ == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if typ != "gauge" && typ != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	err := h.storage.Set(typ, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}
