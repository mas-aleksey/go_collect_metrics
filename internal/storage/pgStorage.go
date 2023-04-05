package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"strings"
)

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
	insertArg := ""

	switch metricIn.MType {
	case "gauge":
		insertArg = fmt.Sprintf("('%s', 'gauge', %v, null)", metricIn.ID, *metricIn.Value)
	case "counter":
		insertArg = fmt.Sprintf("('%s', 'counter', null, %v)", metricIn.ID, *metricIn.Delta)
	}
	stmt := fmt.Sprintf(
		"INSERT INTO metric(name, type, gauge_value, counter_value) VALUES %s %s", insertArg,
		"ON CONFLICT (name, type) DO UPDATE SET gauge_value = excluded.gauge_value, counter_value = metric.counter_value + excluded.counter_value RETURNING name, type, gauge_value, counter_value;",
	)
	row := p.Conn.QueryRow(ctx, stmt)
	err := row.Scan(&metricOut.ID, &metricOut.MType, &metricOut.Value, &metricOut.Delta)
	if err != nil {
		return metricOut, err
	}
	return metricOut, nil
}

func (p *PgStorage) UpdateJSONMetrics(ctx context.Context, metricsIn []utils.JSONMetric) ([]utils.JSONMetric, error) {
	valueStrings := make([]string, 0)
	metricsOut := make([]utils.JSONMetric, 0)

	for _, metric := range metricsIn {
		insertArg := ""
		switch metric.MType {
		case "gauge":
			insertArg = fmt.Sprintf("('%s', 'gauge', %v, null)", metric.ID, *metric.Value)
		case "counter":
			insertArg = fmt.Sprintf("('%s', 'counter', null, %v)", metric.ID, *metric.Delta)
		}
		valueStrings = append(valueStrings, insertArg)
	}
	stmt := fmt.Sprintf("INSERT INTO metric(name, type, gauge_value, counter_value) VALUES %s %s",
		strings.Join(valueStrings, ","),
		"ON CONFLICT (name, type) DO UPDATE SET gauge_value = excluded.gauge_value, counter_value = metric.counter_value + excluded.counter_value RETURNING name, type, gauge_value, counter_value;",
	)
	rows, err := p.Conn.Query(ctx, stmt)
	if err != nil {
		return metricsOut, err
	}
	for rows.Next() {
		metric := utils.JSONMetric{}
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return metricsOut, err
		}
		metricsOut = append(metricsOut, metric)
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
