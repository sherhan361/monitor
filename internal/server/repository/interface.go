package repository

import (
	"github.com/sherhan361/monitor/internal/models"
	"github.com/sherhan361/monitor/internal/server/config"
	"github.com/sherhan361/monitor/internal/server/repository/memory"
)

type Getter interface {
	GetAll() (map[string]float64, map[string]int64)
	Get(typ, name string) (string, error)
	Set(typ, name, value string) error

	GetMetricsByID(id, typ string, key string) (*models.Metric, error)
	SetMetrics(*models.Metric) error

	RestoreMetrics(filename string) error
	WriteMetrics() error
}

func NewMemoryStorage(cfg config.Config) (Getter, error) {
	return memory.New(cfg), nil
}
