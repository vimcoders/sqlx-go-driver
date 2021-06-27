package sqlx

import (
	"context"
	"database/sql"

	"github.com/vimcoders/go-driver"
)

type Conn struct {
	*sql.Conn
}

func (c *Conn) Exec(convert driver.Convertor) (driver.Result, error) {
	return c.ExecContext(context.Background(), convert)
}

func (c *Conn) ExecContext(ctx context.Context, convert driver.Convertor) (driver.Result, error) {
	sql, args := convert.Convert()

	result, err := c.Conn.ExecContext(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Conn) QueryContext(ctx context.Context, scanner driver.Scanner) (err error) {
	sql, args := scanner.Convert()

	rows, err := c.Conn.QueryContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		if err := scanner.Scan(func(dest ...interface{}) error {
			return rows.Scan(dest...)
		}); err != nil {
			return err
		}
	}

	return nil
}

func (c *Conn) Query(scanner driver.Scanner) (err error) {
	return c.QueryContext(context.Background(), scanner)
}

func (c *Conn) Tx(ctx context.Context) (driver.Tx, error) {
	tx, err := c.Conn.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}
