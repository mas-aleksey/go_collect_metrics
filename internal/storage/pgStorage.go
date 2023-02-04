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
	Buffer *Buffer
	Conn   *pgx.Conn
	Config *utils.StorageConfig
}

func (p *PgStorage) Init() error {
	conn, err := pgx.Connect(context.Background(), p.Config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	p.Conn = conn
	err = p.createTable()
	if err != nil {
		return fmt.Errorf("unable to crate table: %v", err)
	}
	err = p.loadInBuffer()
	if err != nil {
		return fmt.Errorf("unable to load data from db: %v", err)
	}
	return nil
}

func (p *PgStorage) Close() {
	err := p.Conn.Close(context.Background())
	if err != nil {
		return
	}
}

func (p *PgStorage) GetConfig() *utils.StorageConfig {
	return p.Config
}

func (p *PgStorage) GetBuffer() *Buffer {
	return p.Buffer
}

func (p *PgStorage) Ping() bool {
	err := p.Conn.Ping(context.Background())
	return err == nil
}

func (p *PgStorage) Save() {
	err := p.flushBuffer()
	if err != nil {
		fmt.Printf("db save error %s", err)
	}
	fmt.Printf("db save success")
}

func (p *PgStorage) SaveIfSyncMode() {
	p.Save()
}

func (p *PgStorage) createTable() error {
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
	returnVal, err := p.Conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	fmt.Printf("value of returnval %s", returnVal)
	return nil
}

func (p *PgStorage) loadInBuffer() error {
	query := "SELECT name, type, gauge_value, counter_value FROM metric"
	rows, err := p.Conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for rows.Next() {
		metric := utils.JSONMetric{}
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		fmt.Println("metric", metric)
		if err != nil {
			log.Fatal(err)
		}
		p.Buffer.PutJSONMetric(metric)
	}
	return nil
}

func (p *PgStorage) flushBuffer() error {
	valueStrings := make([]string, 0)

	p.Buffer.Mutex.RLock()
	defer p.Buffer.Mutex.RUnlock()
	for mName, mValue := range p.Buffer.GaugeMetrics {
		insertArg := fmt.Sprintf("('%s', 'gauge', %v, null)", mName, mValue)
		valueStrings = append(valueStrings, insertArg)
	}
	for mName, mValue := range p.Buffer.CounterMetrics {
		insertArg := fmt.Sprintf("('%s', 'counter', null, %v)", mName, mValue)
		valueStrings = append(valueStrings, insertArg)
	}
	stmt := fmt.Sprintf("INSERT INTO metric(name, type, gauge_value, counter_value) VALUES %s %s",
		strings.Join(valueStrings, ","),
		"ON CONFLICT (name, type) DO UPDATE SET gauge_value = excluded.gauge_value, counter_value = excluded.counter_value;",
	)
	res, err := p.Conn.Exec(context.Background(), stmt)
	if err != nil {
		return err
	}
	fmt.Printf("value of returnval %s", res)
	return nil
}
