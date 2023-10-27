package pg

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(u *url.URL) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, err
	}
	config.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (pg *Postgres) Shutdown() {
	pg.pool.Close()
}

func (pg *Postgres) queryRow(name string, args ...interface{}) (map[string]interface{}, error) {
	q, err := assets.ReadString(name)
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		r := make(map[string]interface{})
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}

		for i, f := range rows.FieldDescriptions() {
			r[string(f.Name)] = vals[i]
		}
		return r, nil
	} else {
		return nil, nil
	}
}

func (pg *Postgres) queryRows(name string, args ...interface{}) ([]map[string]interface{}, error) {
	q, err := assets.ReadString(name)
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rs := make([]map[string]interface{}, 0)
	for rows.Next() {
		r := make(map[string]interface{})
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}

		for i, f := range rows.FieldDescriptions() {
			r[string(f.Name)] = vals[i]
		}
		rs = append(rs, r)
	}

	if len(rs) < 1 {
		return nil, nil
	}
	return rs, nil
}

func (pg *Postgres) queryTable(name string, args ...interface{}) (*model.Table, error) {
	q, err := assets.ReadString(name)
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rs := &model.Table{}
	for _, f := range rows.FieldDescriptions() {
		rs.Columns = append(rs.Columns, string(f.Name))
	}

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rs.Rows = append(rs.Rows, vals)
	}
	return rs, nil
}

func (pg *Postgres) enumList(enum string) ([]string, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	clean_enum, err := conn.Conn().PgConn().EscapeString(enum)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("SELECT ARRAY_AGG(unnest ORDER BY unnest) "+
		"FROM UNNEST(ENUM_RANGE(null::\"%s\")::TEXT[])", clean_enum)
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	if rows.Next() {
		if err := rows.Scan(&values); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (pg *Postgres) enumAddValues(enum string, values ...string) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	e_enum, err := tx.Conn().PgConn().EscapeString(enum)
	if err != nil {
		return err
	}

	for _, value := range values {
		e_value, err := tx.Conn().PgConn().EscapeString(value)
		if err != nil {
			return err
		}
		q := fmt.Sprintf(
			"ALTER TYPE \"%s\" ADD VALUE IF NOT EXISTS '%s'", e_enum, e_value,
		)
		if _, err := tx.Exec(context.Background(), q); err != nil {
			return err
		}
	}
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

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
		return 0, fmt.Errorf("Cannot set value")
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
		return 0, fmt.Errorf("Cannot set value")
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
