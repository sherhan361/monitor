package service

import (
	"encoding/json"
	"github.com/sherhan361/monitor/internal/server/repository"
	"log"
	"os"
	"time"
)

type OutMetrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Backuper struct {
	storeInterval time.Duration
	filename      string
	repo          repository.Getter
}

func NewBackuper(duration time.Duration, filename string, repo repository.Getter) (*Backuper, error) {
	return &Backuper{
		storeInterval: duration,
		filename:      filename,
		repo:          repo,
	}, nil
}

func (b *Backuper) Run() {
	for {
		<-time.After(b.storeInterval)
		err := b.WriteMetric()
		if err != nil {
			log.Println(err)
		}

	}

}

func (b *Backuper) WriteMetric() error {
	var metrics []OutMetrics
	file, err := os.OpenFile(b.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	gauges, counters := b.repo.GetAll()
	for key, value := range gauges {
		var metric = OutMetrics{}
		metric.ID = key
		metric.MType = "gauge"
		metric.Value = value
		metrics = append(metrics, metric)
	}

	for key, value := range counters {
		var metric = OutMetrics{}
		metric.ID = key
		metric.MType = "counter"
		metric.Delta = value
		metrics = append(metrics, metric)
	}

	jsonData, err := json.Marshal(&metrics)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println(err)
	}

	return nil
}
