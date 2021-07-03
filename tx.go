package sqlx

import (
	"context"
	"database/sql"

	"github.com/vimcoders/go-driver"
)

type Tx struct {
	*sql.Tx
}

func (tx *Tx) Exec(list ...interface{}) (driver.Result, error) {
	return tx.Exec(context.Background(), list)
}

func (tx *Tx) ExecContext(ctx context.Context, list ...interface{}) (driver.Result, error) {
	for _, i := range list {
		convert, ok := i.(Convertor)

		if !ok {
			continue
		}

		sql, args := convert.Convert()

		if _, err := tx.Tx.Exec(sql, args...); err == nil {
			continue
		}

		return nil, tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		return nil, tx.Rollback()
	}

	return nil, nil
}
