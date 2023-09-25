package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sherhan361/monitor/internal/common"
	"github.com/sherhan361/monitor/internal/models"
	"github.com/sherhan361/monitor/internal/server/config"
	"log"
	"strconv"
	"time"
)

type DBStor struct {
	config config.Config
	db     *sql.DB
}

func New(cfg config.Config) DBStor {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Println("err:", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	m, err := migrate.New("file://internal/server/repository/db/migrations", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	if _, _, err := m.Version(); err != nil {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}

	return DBStor{
		config: cfg,
		db:     db,
	}
}

func (d DBStor) setGauge(ctx context.Context, key string, newMetricValue float64) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", key, newMetricValue)
	return err
}

func (d DBStor) setCounter(ctx context.Context, key string, newMetricValue int64) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counter.value + $2", key, newMetricValue)
	return err
}

func (d DBStor) Ping() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.db.PingContext(ctx); err != nil {
		panic(err)
	}
	log.Println("ping!")
	return nil
}

func (d DBStor) GetRowsCounter() (*sql.Rows, error) {
	rows, err := d.db.Query("SELECT name, value from counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows, nil
}

func (d DBStor) GetAll() (map[string]float64, map[string]int64) {
	counters := map[string]int64{}
	gauges := map[string]float64{}

	rows, err := d.GetRowsCounter()
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		var name string
		var value int64
		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, nil
		}
		counters[name] = value
	}
	if err = rows.Err(); err != nil {
		return nil, nil
	}

	rows, err = d.db.Query("SELECT name, value from gauge")
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value float64
		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, nil
		}
		gauges[name] = value
	}
	if err = rows.Err(); err != nil {
		return nil, nil
	}

	return gauges, counters
}

func (d DBStor) Get(typ, name string) (string, error) {
	switch typ {
	case "counter":
		var counter string
		row := d.db.QueryRow("SELECT value FROM counter WHERE name = $1", name)
		err := row.Scan(&counter)
		return counter, err
	case "gauge":
		var gauge string
		row := d.db.QueryRow("SELECT value FROM gauge WHERE name = $1", name)
		err := row.Scan(&gauge)
		return gauge, err
	default:
		return "", errors.New("invalid metric type")
	}
}

func (d DBStor) Set(ctx context.Context, typ, name, value string) error {
	switch typ {
	case "counter":
		countValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		err = d.setCounter(ctx, name, countValue)
		if err != nil {
			return err
		}
		return nil
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		err = d.setGauge(ctx, name, floatValue)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid metric type")
	}
}

func (d DBStor) GetMetricsByID(id, typ string, key string) (*models.Metric, error) {
	var input models.Metric

	switch typ {
	case "gauge":
		var gauge float64
		row := d.db.QueryRow("SELECT value FROM gauge WHERE name = $1", id)
		err := row.Scan(&gauge)

		if err == nil {
			input.ID = id
			input.MType = "gauge"
			input.Value = &gauge
			if key != "" {
				input.Hash = common.GetHash(input, key)
			}
		}
	case "counter":
		var counter int64
		row := d.db.QueryRow("SELECT value FROM counter WHERE name = $1", id)
		err := row.Scan(&counter)

		if err == nil {
			input.ID = id
			input.MType = "counter"
			input.Delta = &counter
			if key != "" {
				input.Hash = common.GetHash(input, key)
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

func (d DBStor) SetMetrics(ctx context.Context, metric *models.Metric) error {
	switch metric.MType {
	case "gauge":
		var gauge float64
		if metric.Value == nil {
			gauge = 0
		} else {
			gauge = *metric.Value
		}
		err := d.setGauge(ctx, metric.ID, gauge)
		if err != nil {
			return err
		}

		return nil
	case "counter":
		if metric.Delta == nil {
			return errors.New("invalid params")
		}
		err := d.setCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("invalid metric type")
	}
}

func (d DBStor) SetMetricsBatch(ctx context.Context, MetricBatch []models.Metric) error {
	MetricValueBatch := models.Metric{}
	for _, OneMetric := range MetricBatch {
		MetricValueBatch = models.Metric{
			ID:    OneMetric.ID,
			MType: OneMetric.MType,
			Delta: OneMetric.Delta,
			Value: OneMetric.Value,
		}
		err := d.SetMetrics(ctx, &MetricValueBatch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DBStor) RestoreMetrics(filename string) error {
	return nil
}

func (d DBStor) WriteMetrics() error {
	return nil
}
