package repository

import (
	"github.com/sherhan361/monitor/internal/server/repository/memory"
)

type Getter interface {
	GetAll() (map[string]float64, map[string]int64)
	Get(typ, name string) (string, error)
	Set(typ, name, value string) error

	GetMetricsByID(id, typ string) (*memory.Metrics, error)
	SetMetrics(*memory.Metrics) error
}

func NewGetter() (Getter, error) {
	return memory.New(), nil
}
