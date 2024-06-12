package pg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

func structFieldMatcher(fieldName, ch string) bool {
	ch = strings.ReplaceAll(ch, "_", "")
	return strings.EqualFold(fieldName, ch)
}

func rowsToStruct(rows pgx.Rows, a interface{}) (int, error) {
	rv := reflect.ValueOf(a)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.CanSet() {
		return 0, fmt.Errorf("cannot set value")
	}
	switch rv.Kind() {
	case reflect.Slice:
		for rowCount := 0; ; {
			var elem reflect.Value
			switch typ := rv.Type().Elem(); typ.Kind() {
			case reflect.Ptr:
				elem = reflect.New(typ.Elem())
			case reflect.Struct:
				elem = reflect.New(typ).Elem()
			}
			if c, err := rowsToStruct(rows, elem.Addr().Interface()); c > 0 {
				if err != nil {
					return 0, err
				}
				rowCount += c
				rv.Set(reflect.Append(rv, elem))
				continue
			}
			if rowCount == 1 {
				return rowToStruct(rows, a)
			}
			return rowCount, nil
		}
	default:
		if rows.Next() {
			return rowToStruct(rows, a)
		}
		defer rows.Close()
	}
	return 0, nil
}

func rowToStruct(rows pgx.Rows, a interface{}) (int, error) {
	rv := reflect.ValueOf(a)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.CanSet() {
		return 0, fmt.Errorf("cannot set value")
	}
	switch rv.Kind() {
	case reflect.Slice:
		var elem reflect.Value
		switch typ := rv.Type().Elem(); typ.Kind() {
		case reflect.Ptr:
			elem = reflect.New(typ.Elem())
		case reflect.Struct:
			elem = reflect.New(typ).Elem()
		}
		for rowCount := 0; ; {
			if c, err := rowToStruct(rows, elem.Addr().Interface()); c > 0 {
				if err != nil {
					return 0, err
				}
				rowCount += c
				rv.Set(reflect.Append(rv, elem))
				continue
			}
			return rowCount, nil
		}
	case reflect.Struct:
		var columnNames [][]string
		cols := rows.FieldDescriptions()
		for _, c := range cols {
			parts := strings.Split(string(c.Name), ".")
			columnNames = append(columnNames, parts)
		}
		for i := 0; i < rv.NumField(); i++ {
			if rv.Field(i).Kind() == reflect.Ptr {
				if rv.Field(i).Type().Elem().Kind() == reflect.Struct {
					rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
				}
			}
		}
		ptrsFields := make([]interface{}, len(cols))
		values, err := rows.Values()
		if err != nil {
			return 0, err
		}
		if values == nil {
			return 0, fmt.Errorf("rows.Values() returned nil")
		}
		for i, c := range columnNames {
			if values[i] == nil {
				ptrsFields[i] = nil
				continue
			}
			f := rv
			for _, ch := range c {
				if f.Kind() == reflect.Ptr {
					f = f.Elem()
				}
				f = f.FieldByNameFunc(func(fieldName string) bool {
					return structFieldMatcher(fieldName, string(ch))
				})
				if !f.IsValid() {
					break
				}
				if f.IsValid() {
					for j := 0; j < len(cols); j++ {
						fieldName := string(cols[j].Name)
						if strings.Contains(fieldName, ".") {
							parts := strings.Split(fieldName, ".")
							fieldName = parts[len(parts)-1]
						}
						if strings.EqualFold(string(fieldName), string(ch)) {
							ptrsFields[i] = f.Addr().Interface()
						}
					}
				}
			}
		}
		err = rows.Scan(ptrsFields...)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return 0, nil
}
