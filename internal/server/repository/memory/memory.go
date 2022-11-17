package memory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

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
	return m.Gauges, m.Counters
}

func (m *MemStorage) Get(typ, name string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value := ""
	if typ == "gauge" {
		val, ok := m.Gauges[name]
		if ok {
			value = fmt.Sprintf("%.3f", val)
		}
	}
	if typ == "counter" {
		val, ok := m.Counters[name]
		if ok {
			value = fmt.Sprintf("%d", val)
		}
	}
	if value == "" {
		return "", errors.New("metric not found")
	}
	return value, nil
}

func (m *MemStorage) Set(typ, name, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if typ == "counter" {
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
	}

	if typ == "gauge" {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		m.Gauges[name] = floatValue
		return nil
	}

	return errors.New("invalid request params")
}
