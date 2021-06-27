package sqlx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vimcoders/go-driver"
)

type Queryer struct {
	Where   interface{}
	Scanner interface{}
}

func (q *Queryer) scanner() interface{} {
	t := reflect.TypeOf(q.Scanner).Elem()

	if t == nil {
		return nil
	}

	switch t.Kind() {
	case reflect.Slice:
		t = t.Elem()

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		scanner := reflect.New(t).Interface()
		v := reflect.ValueOf(q.Scanner).Elem()
		v.Set(reflect.Append(v, reflect.ValueOf(scanner)))

		return scanner
	case reflect.Struct:
		return q.Scanner
	}

	return nil
}

func (q *Queryer) TableName() string {
	if table, ok := q.Scanner.(driver.Table); ok && len(table.TableName()) > 0 {
		return table.TableName()
	}

	t := reflect.TypeOf(q.Scanner).Elem()

	switch t.Kind() {
	case reflect.Struct:
		return strings.ToLower(t.Name())
	case reflect.Slice:
		t = t.Elem()

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return strings.ToLower(t.Name())
	}

	return ""
}

func (q *Queryer) where() (query string, args []interface{}) {
	t, v := reflect.TypeOf(q.Where).Elem(), reflect.ValueOf(q.Where).Elem()

	if t == nil || t.Kind() != reflect.Struct {
		return "", nil
	}

	var values []string

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
		default:
			continue
		}

		colName := t.Field(i).Tag.Get("db")

		if len(colName) <= 0 {
			continue
		}

		values = append(values, fmt.Sprintf("`%v`=?", colName))
		args = append(args, v.Field(i).Interface())
	}

	return strings.Join(values, " AND "), args
}

func (q *Queryer) Convert() (query string, args []interface{}) {
	t := reflect.TypeOf(q.Scanner).Elem()

	if t == nil {
		return "", nil
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	var values []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		switch field.Type.Kind() {
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
		case reflect.Array:
		case reflect.Slice:
		default:
			continue
		}

		values = append(values, fmt.Sprintf("`%v`", field.Tag.Get("db")))
	}

	if q.Where == nil {
		return fmt.Sprintf("SELECT %v FROM `%v`;", strings.Join(values, ","), q.TableName()), nil
	}

	query, args = q.where()

	return fmt.Sprintf("SELECT %v FROM `%v` WHERE %v;", strings.Join(values, ","), q.TableName(), query), args
}

func (q *Queryer) pointer(value interface{}) (dest []interface{}) {
	v := reflect.ValueOf(value).Elem()

	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
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
		default:
			dest = append(dest, new(string))
			continue
		}

		dest = append(dest, v.Field(i).Addr().Interface())
	}

	return dest
}

func (q *Queryer) Scan(scan func(dest ...interface{}) error) error {
	scanner := q.scanner()

	dest := q.pointer(scanner)

	if err := scan(dest...); err != nil {
		return err
	}

	v := reflect.ValueOf(scanner).Elem()

	for i := 0; i < len(dest); i++ {
		switch v.Field(i).Kind() {
		case reflect.Map:
		case reflect.Slice:
		case reflect.Array:
		default:
			continue
		}

		str, ok := dest[i].(*string)

		if !ok {
			continue
		}

		decoder, ok := v.Field(i).Addr().Interface().(driver.Unmarshaler)

		if !ok {
			continue
		}

		if err := decoder.Unmarshal(*str); err != nil {
			return err
		}
	}

	return nil
}

func WithQueryer(where, scanner interface{}) driver.Convertor {
	if scanner == nil {
		return nil
	}

	switch reflect.TypeOf(where).Kind() {
	case reflect.Ptr:
	default:
		return nil
	}

	t := reflect.TypeOf(scanner)

	if t == nil || t.Kind() != reflect.Ptr {
		return nil
	}

	t = t.Elem()

	switch t.Kind() {
	case reflect.Slice:
		if t.Elem().Kind() != reflect.Ptr {
			return nil
		}
	case reflect.Struct:
	default:
		return nil
	}

	return &Queryer{
		Where:   where,
		Scanner: scanner,
	}
}
