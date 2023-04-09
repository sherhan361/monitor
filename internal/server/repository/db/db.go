package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sherhan361/monitor/internal/models"
	"github.com/sherhan361/monitor/internal/server/config"
	"strconv"
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
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS counter (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	if err != nil {
		fmt.Println("create counter table error:", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS gauge (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		fmt.Println("create gauge table error:", err)
	}

	return DBStor{
		config: cfg,
		db:     db,
	}
}

func (d DBStor) setGauge(key string, newMetricValue float64) error {
	_, err := d.db.Exec("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", key, newMetricValue)
	return err
}

func (d DBStor) setCounter(key string, newMetricValue int64) error {
	_, err := d.db.Exec("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counter.value + $2", key, newMetricValue)
	return err
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
	switch typ {
	case "counter":
		countValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		err = d.setCounter(name, countValue)
		if err != nil {
			return err
		}
		return nil
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		err = d.setGauge(name, floatValue)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid metric type")
	}
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
