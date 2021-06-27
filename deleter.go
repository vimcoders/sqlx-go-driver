package sqlx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vimcoders/go-driver"
)

type Deleter struct {
	Where   interface{}
	Deleter interface{}
}

func (d *Deleter) TableName() string {
	if table, ok := d.Deleter.(driver.Table); ok {
		return table.TableName()
	}

	if t := reflect.TypeOf(d.Deleter).Elem(); t != nil {
		return strings.ToLower(t.Name())
	}

	return ""
}

func (d *Deleter) Convert() (sql string, args []interface{}) {
	t, v := reflect.TypeOf(d.Where).Elem(), reflect.ValueOf(d.Where).Elem()

	if t == nil {
		return
	}

	var keys []string

	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int:
		case reflect.Int8:
		case reflect.Int16:
		case reflect.Int32:
		case reflect.Int64:
		case reflect.Uint:
		case reflect.Uint8:
		case reflect.Uint16:
		case reflect.Uint32:
		case reflect.Uint64:
		case reflect.String:
		default:
			continue
		}

		colName := t.Field(i).Tag.Get("db")

		if len(colName) <= 0 {
			continue
		}

		args = append(args, v.Field(i).Interface())
		keys = append(keys, fmt.Sprintf("`%v`=?", colName))
	}

	return fmt.Sprintf("DELETE FROM `%v` WHERE %v;", d.TableName(), strings.Join(keys, " AND ")), args
}

func (d *Deleter) Scan(scan func(dest ...interface{}) error) error {
	return nil
}

func WithDeleter(where, deleter interface{}) driver.Convertor {
	switch reflect.TypeOf(where).Kind() {
	case reflect.Ptr:
	default:
		return nil
	}

	return &Deleter{
		Where:   where,
		Deleter: deleter,
	}
}
