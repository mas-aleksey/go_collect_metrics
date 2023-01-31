package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"os"
)

type PgStorage struct {
	Buffer *Buffer
	Dns    string
	Conn   *pgx.Conn
	Config *utils.StorageConfig
}

func (p *PgStorage) Init() error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	p.Conn = conn
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
	if err != nil {
		return false
	}
	return true
}

func (p *PgStorage) Save() {
	//TODO implement me
	panic("implement me")
}

func (p *PgStorage) SaveIfSyncMode() {
	//TODO implement me
	panic("implement me")
}
