package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sherhan361/monitor/internal/models"
	"github.com/sherhan361/monitor/internal/server/config"
	"time"
)

type DBStor struct {
	config config.Config
	db     *sql.DB
}

func New(cfg config.Config) DBStor {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		fmt.Println("err:", err)
	}
	return DBStor{
		config: cfg,
		db:     db,
	}
}

func (d DBStor) Ping() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.db.PingContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("ping!")
	return nil
}

func (d DBStor) GetAll() (map[string]float64, map[string]int64) {
	return nil, nil
}

func (d DBStor) Get(typ, name string) (string, error) {
	return "", nil
}

func (d DBStor) Set(typ, name, value string) error {
	return nil
}

func (d DBStor) GetMetricsByID(id, typ string, key string) (*models.Metric, error) {
	return nil, nil
}

func (d DBStor) SetMetrics(metric *models.Metric) error {
	return nil
}

func (d DBStor) RestoreMetrics(filename string) error {
	return nil
}

func (d DBStor) WriteMetrics() error {
	return nil
}
