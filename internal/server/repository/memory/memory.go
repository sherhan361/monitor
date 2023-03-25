package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sherhan361/monitor/internal/common"
	"github.com/sherhan361/monitor/internal/models"
	"github.com/sherhan361/monitor/internal/server/config"
	"log"
	"os"
	"strconv"
	"sync"
)

type MemStorage struct {
	mutex    *sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
	Config   config.Config
}

func New(cfg config.Config) *MemStorage {
	return &MemStorage{
		mutex:    &sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
		Config:   cfg,
	}
}

func (m *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	gauges := make(map[string]float64, len(m.Gauges))
	for key, value := range m.Gauges {
		gauges[key] = value
	}
	counters := make(map[string]int64, len(m.Counters))
	for key, value := range m.Counters {
		counters[key] = value
	}
	return gauges, counters
}

func (m *MemStorage) Get(typ, name string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value := ""
	switch typ {
	case "gauge":
		val, ok := m.Gauges[name]
		if ok {
			value = fmt.Sprintf("%.3f", val)
		} else {
			return "", errors.New("metric name not found in Gauges")
		}

	case "counter":
		val, ok := m.Counters[name]
		if ok {
			value = fmt.Sprintf("%d", val)
		} else {
			return "", errors.New("metric name not found in Counters")
		}
	default:
		return "", errors.New("metric type not found")
	}
	return value, nil
}

func (m *MemStorage) Set(typ, name, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch typ {
	case "counter":
		countValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		_, ok := m.Counters[name]
		if ok {
			m.Counters[name] = m.Counters[name] + countValue
		} else {
			m.Counters[name] = countValue
		}
		return nil
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		m.Gauges[name] = floatValue
		return nil
	default:
		return errors.New("invalid metric type")
	}
}

func (m *MemStorage) GetMetricsByID(id, typ string, signKey string) (*models.Metric, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var input models.Metric

	switch typ {
	case "gauge":
		v, ok := m.Gauges[id]
		if ok {
			input.ID = id
			input.MType = "gauge"
			input.Value = &v
			if signKey != "" {
				input.Hash = common.GetHash(input, signKey)
			}
		}
	case "counter":
		v, ok := m.Counters[id]
		if ok {
			input.ID = id
			input.MType = "counter"
			input.Delta = &v
			if signKey != "" {
				input.Hash = common.GetHash(input, signKey)
			}
		}
	default:
		return nil, errors.New("invalid metric type")
	}

	if input.ID == "" {
		return nil, errors.New("not found")
	}

	return &input, nil
}

func (m *MemStorage) SetMetrics(metrics *models.Metric) error {
	m.mutex.Lock()

	switch metrics.MType {
	case "gauge":
		if metrics.Value == nil {
			m.Gauges[metrics.ID] = 0
		} else {
			m.Gauges[metrics.ID] = *metrics.Value
		}
		m.mutex.Unlock()
		if m.Config.StoreInterval == 0 {
			err := m.WriteMetrics()
			if err != nil {
				return err
			}
		}
		return nil
	case "counter":
		if metrics.Delta == nil {
			return errors.New("invalid params")
		}
		value, ok := m.Counters[metrics.ID]

		if ok {
			m.Counters[metrics.ID] = value + *metrics.Delta
		} else {
			m.Counters[metrics.ID] = *metrics.Delta
		}
		m.mutex.Unlock()
		if m.Config.StoreInterval == 0 {
			err := m.WriteMetrics()
			if err != nil {
				return err
			}
		}
		return nil
	default:
		m.mutex.Unlock()
		return errors.New("invalid metric type")
	}
}

func (m *MemStorage) RestoreMetrics(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var metrics []models.Metric
	err = json.Unmarshal([]byte(content), &metrics)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			if metric.Value != nil {
				m.Gauges[metric.ID] = *metric.Value
			}
		case "counter":
			if metric.Delta != nil {
				m.Counters[metric.ID] = *metric.Delta
			}
		}
	}

	return nil
}

func (m *MemStorage) WriteMetrics() error {
	var metrics []models.WriteMetric
	file, err := os.OpenFile(m.Config.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	gauges, counters := m.GetAll()
	for key, value := range gauges {
		var metric = models.WriteMetric{}
		metric.ID = key
		metric.MType = "gauge"
		metric.Value = value
		metrics = append(metrics, metric)
	}

	for key, value := range counters {
		var metric = models.WriteMetric{}
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
