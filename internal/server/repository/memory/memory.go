package memory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MemStorage struct {
	mutex    *sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		mutex:    &sync.RWMutex{},
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
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
		fmt.Println("m.Counters[PollCount]:", m.Counters["PollCount"])
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

func (m *MemStorage) GetMetricsByID(id, typ string) (*Metrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var input Metrics

	switch typ {
	case "gauge":
		v, ok := m.Gauges[id]
		if ok {
			input.ID = id
			input.MType = "gauge"
			input.Value = &v
		}
	case "counter":
		v, ok := m.Counters[id]
		if ok {
			input.ID = id
			input.MType = "counter"
			input.Delta = &v
		}
	default:
		return nil, errors.New("invalid metric type")
	}

	if input.ID == "" {
		return nil, errors.New("not found")
	}

	return &input, nil
}

func (m *MemStorage) SetMetrics(metrics *Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch metrics.MType {
	case "gauge":
		if metrics.Value == nil {
			m.Gauges[metrics.ID] = 0
		} else {
			m.Gauges[metrics.ID] = *metrics.Value
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
		return nil
	default:
		return errors.New("invalid metric type")
	}
}
