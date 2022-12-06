package pg

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gmc/assets"
	"gmc/db/model"
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

func rowToStruct(r pgx.Rows, a interface{}) int {
	rv := reflect.ValueOf(a)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice:
		var elem reflect.Value
		typ := rv.Type().Elem()

		if typ.Kind() == reflect.Ptr {
			elem = reflect.New(typ.Elem())
		}
		if typ.Kind() == reflect.Struct {
			elem = reflect.New(typ).Elem()
		}
		rowCount := 0
		for {
			if rCount := rowToStruct(r, elem.Addr().Interface()); rCount > 0 {
				rowCount = rowCount + rCount
				rv.Set(reflect.Append(rv, elem))
			} else {
				break
			}
		}
		return rowCount
	case reflect.Struct:
		if r.Next() {
			fieldDescriptions := r.FieldDescriptions()
			var columns []string
			for _, col := range fieldDescriptions {
				columns = append(columns, string(col.Name))
			}
			columnValues, _ := r.Values()
			for j, val := range columnValues {
				switch val.(type) {
				case pgtype.TextArray:
					s := val.(pgtype.TextArray)
					var s_arr []string
					s.AssignTo(&s_arr)
					val = s_arr
				}
				for i := 0; i < rv.NumField(); i++ {
					if reflect.TypeOf(val) != rv.Field(i).Type() {
						continue
					}
					if !strings.EqualFold(columns[j], rv.Type().Field(i).Name) {
						continue
					}
					if rv.CanSet() {
						rv.FieldByName(rv.Type().Field(i).Name).Set(reflect.ValueOf(val))
					}
				}
			}
			return 1
		}
	}
	return 0
}
