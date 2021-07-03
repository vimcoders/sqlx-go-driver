package sqlx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vimcoders/go-driver"
)

type Config struct {
	DriverName string
	Usr        string
	Pwd        string
	Addr       string
	DB         string
}

func Connect(cfg *Config) (driver.Connector, error) {
	//TODO::decode pwd

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True", cfg.Usr, cfg.Pwd, cfg.Addr, cfg.DB)

	connect, err := sql.Open(cfg.DriverName, dsn)

	if err != nil {
		return nil, err
	}

	return &Connector{
		DB: connect,
	}, nil
}

type Connector struct {
	*sql.DB
}

func (c *Connector) Tx(ctx context.Context) (driver.Execer, error) {
	tx, err := c.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}
