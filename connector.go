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

func (c *Connector) Conn(ctx context.Context) (driver.Conn, error) {
	return nil, nil
	//dc, err := c.db.Conn(ctx)

	//if err != nil {
	//	return nil, err
	//}

	//return &connect{
	//	conn: dc,
	//}, nil
}

func (c *Connector) Tx(ctx context.Context) (driver.Tx, error) {
	return nil, nil
	//tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})

	//if err != nil {
	//	return nil, err
	//}

	//return &trans{tx}, nil
}

func (c *Connector) SetMaxOpenConns(n int) {
	//c.db.SetMaxOpenConns(n)
}

func (c *Connector) Close() (err error) {
	return nil
	//return c.db.Close()
}
