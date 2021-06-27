package sqlx

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vimcoders/go-driver"
)

type Inserter struct {
	Where    interface{}
	Inserter interface{}
}

func (i *Inserter) TableName() string {
	if table, ok := i.Inserter.(driver.Table); ok {
		return table.TableName()
	}

	if t := reflect.TypeOf(i.Inserter).Elem(); t != nil {
		return strings.ToLower(t.Name())
	}

	return ""
}

func (i *Inserter) where() (keys []string, args []interface{}) {
	if i.Where == nil {
		return
	}

	t, v := reflect.TypeOf(i.Where).Elem(), reflect.ValueOf(i.Where).Elem()

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
		default:
			continue
		}

		colName := t.Field(i).Tag.Get("db")

		if len(colName) <= 0 {
			continue
		}

		keys = append(keys, fmt.Sprintf("`%v`", colName))
		args = append(args, v.Field(i).Interface())
	}

	return keys, args
}

func (i *Inserter) Convert() (query string, args []interface{}) {
	t, v := reflect.TypeOf(i.Inserter).Elem(), reflect.ValueOf(i.Inserter).Elem()

	if t == nil {
		return
	}

	var keys []string

	keys, args = i.where()

	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Bool:
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
		case reflect.Struct:
		case reflect.Slice:
		default:
			continue
		}

		colName := t.Field(i).Tag.Get("db")

		if len(colName) <= 0 {
			continue
		}

		keys = append(keys, fmt.Sprintf("`%v`", colName))

		if encoder, ok := v.Field(i).Addr().Interface().(driver.Marshaler); ok {
			str, err := encoder.Marshal()

			if err != nil {
				return
			}

			args = append(args, str)
			continue
		}

		if t, ok := v.Field(i).Interface().(time.Time); ok {
			args = append(args, t.Format("2006-01-02 15:04:05"))
			continue
		}

		args = append(args, v.Field(i).Interface())
	}

	values := make([]string, len(args))

	for i := 0; i < len(args); i++ {
		values[i] = "?"
	}

	return fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v);", i.TableName(), strings.Join(keys, ","), strings.Join(values, ",")), args
}

func (i *Inserter) Scan(scan func(dest ...interface{}) error) error {
	return nil
}

func WithInserter(where, inserter interface{}) driver.Convertor {
	if inserter == nil {
		return nil
	}

	switch reflect.TypeOf(inserter).Kind() {
	case reflect.Ptr:
	default:
		return nil
	}

	if where != nil {
		switch reflect.TypeOf(where).Kind() {
		case reflect.Ptr:
		default:
			return nil
		}
	}

	return &Inserter{
		Where:    where,
		Inserter: inserter,
	}
}
