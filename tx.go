package sqlx

import (
	"context"
	"database/sql"

	"github.com/vimcoders/go-driver"
)

type Tx struct {
	*sql.Tx
}

func (tx *Tx) Exec(convert driver.Convertor) (driver.Result, error) {
	return tx.ExecContext(context.Background(), convert)
}

func (tx *Tx) ExecContext(ctx context.Context, convert driver.Convertor) (driver.Result, error) {
	sql, args := convert.Convert()

	result, err := tx.Tx.ExecContext(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (tx *Tx) Query(scanner driver.Scanner) (err error) {
	return tx.QueryContext(context.Background(), scanner)
}

func (tx *Tx) QueryContext(ctx context.Context, scanner driver.Scanner) (err error) {
	sql, args := scanner.Convert()

	rows, err := tx.Tx.QueryContext(ctx, sql, args...)

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
