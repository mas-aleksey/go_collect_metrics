package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
)

var counterStmt = `INSERT INTO metric(name, type, gauge_value, counter_value) VALUES ($1, 'counter', null, $2) 
		ON CONFLICT (name, type) DO UPDATE SET 
		    counter_value = metric.counter_value + excluded.counter_value 
		RETURNING name, type, gauge_value, counter_value;`

var gaugeStmt = `INSERT INTO metric(name, type, gauge_value, counter_value) VALUES ($1, 'gauge', $2, null) 
		ON CONFLICT (name, type) DO UPDATE SET 
		    gauge_value = excluded.gauge_value
		RETURNING name, type, gauge_value, counter_value;`

// PgStorage - структура для работы с бд Postgres
type PgStorage struct {
	Conn   *pgx.Conn
	Config *utils.StorageConfig
}

func (p *PgStorage) Init(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, p.Config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	p.Conn = conn
	err = p.createTable(ctx)
	if err != nil {
		return fmt.Errorf("unable to crate table: %v", err)
	}
	return nil
}

func (p *PgStorage) Close(ctx context.Context) {
	err := p.Conn.Close(ctx)
	if err != nil {
		return
	}
}

func (p *PgStorage) Ping(ctx context.Context) bool {
	err := p.Conn.Ping(ctx)
	return err == nil
}

func (p *PgStorage) UpdateJSONMetric(ctx context.Context, metricIn utils.JSONMetric) (utils.JSONMetric, error) {
	metricOut := utils.JSONMetric{}
	var row pgx.Row
	switch metricIn.MType {
	case "counter":
		row = p.Conn.QueryRow(ctx, counterStmt, metricIn.ID, *metricIn.Delta)
	case "gauge":
		row = p.Conn.QueryRow(ctx, gaugeStmt, metricIn.ID, *metricIn.Value)
	}
	err := row.Scan(&metricOut.ID, &metricOut.MType, &metricOut.Value, &metricOut.Delta)
	if err != nil {
		return metricOut, err
	}
	return metricOut, nil
}

func (p *PgStorage) UpdateJSONMetrics(ctx context.Context, metricsIn []utils.JSONMetric) ([]utils.JSONMetric, error) {
	metricsOut := make([]utils.JSONMetric, 0)
	tx, err := p.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return metricsOut, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	for _, metric := range metricsIn {
		out, err := p.UpdateJSONMetric(ctx, metric)
		if err != nil {
			return metricsOut, err
		}
		metricsOut = append(metricsOut, out)
	}
	return metricsOut, nil
}

func (p *PgStorage) GetJSONMetric(ctx context.Context, mName, mType string) (utils.JSONMetric, error) {
	metric := utils.JSONMetric{}
	query := fmt.Sprintf("SELECT name, type, gauge_value, counter_value FROM metric WHERE name='%s' and type='%s';", mName, mType)
	row := p.Conn.QueryRow(ctx, query)
	err := row.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
	if err != nil {
		return metric, err
	}
	return metric, nil
}

func (p *PgStorage) GetAllMetrics(ctx context.Context) ([]utils.JSONMetric, error) {
	metrics := make([]utils.JSONMetric, 0)
	query := "SELECT name, type, gauge_value, counter_value FROM metric;"
	rows, err := p.Conn.Query(ctx, query)
	if err != nil {
		return metrics, err
	}
	for rows.Next() {
		metric := utils.JSONMetric{}
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (p *PgStorage) createTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS metric(
			id SERIAL PRIMARY KEY,
			name varchar(45) NOT NULL,
			type varchar(15) NOT NULL,
			gauge_value double precision,
			counter_value bigint,
		    UNIQUE (name, type)
		);
	`
	returnVal, err := p.Conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	log.Printf("create table: %s", returnVal)
	return nil
}
