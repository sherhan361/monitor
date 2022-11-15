package memory

import (
	"errors"
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
