package sqlx

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vimcoders/go-driver"
)

type Updater struct {
	Where   interface{}
	Updater interface{}
}

func (u *Updater) TableName() string {
	if u.Updater == nil {
		return ""
	}

	if table, ok := u.Updater.(driver.Table); ok && len(table.TableName()) > 0 {
		return table.TableName()
	}

	if t := reflect.TypeOf(u.Updater).Elem(); t != nil {
		return strings.ToLower(t.Name())
	}

	return ""
}

func (u *Updater) where() (query string, args []interface{}) {
	if u.Where == nil {
		return "", nil
	}

	t, v := reflect.TypeOf(u.Where).Elem(), reflect.ValueOf(u.Where).Elem()

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

func (u *Updater) Convert() (sql string, args []interface{}) {
	t, v := reflect.TypeOf(u.Updater).Elem(), reflect.ValueOf(u.Updater).Elem()

	if t == nil {
		return
	}

	var values []string

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
		case reflect.Struct:
		case reflect.Slice:
		default:
			continue
		}

		colName := t.Field(i).Tag.Get("db")

		if len(colName) <= 0 {
			continue
		}

		values = append(values, fmt.Sprintf("`%v`=?", colName))

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

	query, whereArgs := u.where()

	args = append(args, whereArgs...)

	return fmt.Sprintf("UPDATE `%v` SET %v WHERE %v;", u.TableName(), strings.Join(values, ","), query), args
}

func (u *Updater) Scan(scan func(dest ...interface{}) error) error {
	return nil
}

func WithUpdater(where, updater interface{}) driver.Convertor {
	switch reflect.TypeOf(updater).Kind() {
	case reflect.Ptr:
	default:
		return nil
	}

	switch reflect.TypeOf(where).Kind() {
	case reflect.Ptr:
	default:
		return nil
	}

	return &Updater{
		Where:   where,
		Updater: updater,
	}
}
