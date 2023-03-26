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
	Config config.Config
	db     *sql.DB
}

func (D DBStor) Ping() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := D.db.PingContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("ping!")
	return nil
}

func (D DBStor) GetAll() (map[string]float64, map[string]int64) {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) Get(typ, name string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) Set(typ, name, value string) error {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) GetMetricsByID(id, typ string, key string) (*models.Metric, error) {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) SetMetrics(metric *models.Metric) error {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) RestoreMetrics(filename string) error {
	//TODO implement me
	panic("implement me")
}

func (D DBStor) WriteMetrics() error {
	//TODO implement me
	panic("implement me")
}

func New(cfg config.Config) *DBStor {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		fmt.Println("err:", err)
	}
	defer db.Close()
	return &DBStor{
		Config: cfg,
		db:     db,
	}
}
