package pg

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgtype"
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

func rowToStruct(r pgx.Rows, a interface{}) int {
	rv := reflect.ValueOf(a)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.CanSet() {
		return 0
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
			if c := rowToStruct(r, elem.Addr().Interface()); c > 0 {
				rowCount += c
				rv.Set(reflect.Append(rv, elem))
				continue
			}
			return rowCount
		}
	case reflect.Struct:
		if r.Next() {
			columnValues, _ := r.Values()
			var slices [][]string
			for i, val := range columnValues {
				fieldName := string(r.FieldDescriptions()[i].Name)
				parts := strings.Split(fieldName, ".")
				slice_element := []string{}
				for j := 0; j < len(parts); j++ {
					parts[j] = strings.Replace(parts[j], "_", "", -1)
					slice_element = append(slice_element, parts[j])
				}
				slices = append(slices, slice_element)

				for j := 0; j < rv.NumField(); j++ {
					if rv.Field(j).Kind() == reflect.Slice {
						if len(slice_element) == 1 && strings.EqualFold(rv.Type().Field(j).Name, slice_element[0]) {
							switch val.(type) {
							case []interface{}:
								rec, ok := val.([]interface{})
								if !ok {
									return 0
								}
								var elem reflect.Value
								for i := 0; i < len(rec); i++ {
									switch typ := rv.Field(j).Type().Elem(); typ.Kind() {
									case reflect.Ptr:
										elem = reflect.New(typ.Elem())
									case reflect.Struct:
										elem = reflect.New(typ).Elem()
									}
									switch rec[i].(type) {
									case map[string]interface{}:
										for k, v := range rec[i].(map[string]interface{}) {
											if v != nil {
												switch elem.FieldByName(k).Kind() {
												case reflect.Int32:
													elem.FieldByName(k).Set(reflect.ValueOf(int32(v.(float64))))
												default:
													elem.FieldByName(k).Set(reflect.ValueOf(v))
												}
											}
										}
									}
									rv.Field(j).Set(reflect.Append(rv.Field(j), elem))
								}
							case pgtype.TextArray:
								s := val.(pgtype.TextArray)
								var s_arr []string
								s.AssignTo(&s_arr)
								if reflect.TypeOf(s_arr) == rv.Field(j).Type() {
									rv.Field(j).Set(reflect.ValueOf(s_arr))
								}
							}

						}
						if !strings.EqualFold(rv.Type().Name(), slice_element[0]) {
							continue
						}
						if strings.EqualFold(rv.Type().Field(j).Name, slice_element[0]) {
							var elem reflect.Value
							switch typ := rv.Field(j).Type().Elem(); typ.Kind() {
							case reflect.Ptr:
								elem = reflect.New(typ.Elem())
							case reflect.Struct:
								elem = reflect.New(typ).Elem()
							}
							if elem.Kind() == reflect.Struct {
								for l := 0; l < elem.NumField(); l++ {
									if reflect.TypeOf(val) == elem.Field(l).Type() {
										elem.Field(l).Set(reflect.ValueOf(val))
									}
								}
								rv.Field(j).Set(reflect.Append(rv.Field(j), elem))
							}
						}
					}
					if rv.Field(j).Kind() == reflect.Struct {
						if strings.EqualFold(rv.Type().Field(j).Name, slice_element[0]) {
							for k := 0; k < rv.Field(j).NumField(); k++ {
								if len(slice_element) >= 2 {
									if strings.EqualFold(rv.Field(j).Type().Field(k).Name, slice_element[1]) {
										switch val.(type) {
										case int16, int32:
											if rv.Field(j).Field(k).Kind() == reflect.Ptr {
												p := reflect.New(reflect.TypeOf(val))
												p.Elem().Set(reflect.ValueOf(val))
												rv.Field(j).Field(k).Set(p)
											} else {
												rv.Field(j).Field(k).Set(reflect.ValueOf(val))
											}
										case pgtype.TextArray:
											s := val.(pgtype.TextArray)
											var s_arr []string
											s.AssignTo(&s_arr)
											if reflect.TypeOf(s_arr) == rv.Field(j).Field(k).Type() {
												rv.Field(j).Field(k).Set(reflect.ValueOf(s_arr))
											}
										case pgtype.Numeric:
											n := val.(pgtype.Numeric)
											var nf float64
											n.AssignTo(&nf)
											switch rv.Field(j).Field(k).Kind() {
											case reflect.Ptr:
												if reflect.TypeOf(nf) == rv.Field(j).Field(k).Type() {
													rv.Field(j).Set(reflect.ValueOf(nf))
												} else if rv.Field(j).Field(k).Type().Elem() == reflect.TypeOf(nf) {
													rv.Field(j).Field(k).Set(reflect.ValueOf(&nf))
												}
											case reflect.Struct:
												for l := 0; l < rv.Field(j).Field(k).Field(l).NumField(); l++ {
													if reflect.TypeOf(nf) == rv.Field(j).Field(k).Field(l).Type() {
														rv.Field(j).Field(k).Field(l).Set(reflect.ValueOf(nf))
													}
												}
											default:
												if reflect.TypeOf(nf) == rv.Field(j).Field(k).Type() {
													rv.Field(j).Field(k).Set(reflect.ValueOf(nf))
												}
											}
										case time.Time:
											t, ok := val.(time.Time)
											if ok {
												rv.Field(j).Field(k).Set(reflect.ValueOf(t))
											}
										default:
											if !strings.EqualFold(rv.Type().Field(j).Name, slice_element[0]) {
												continue
											}
											if !strings.EqualFold(rv.Field(j).Type().Field(k).Name, slice_element[1]) {
												continue
											}
											if rv.Field(j).Field(k).Kind() == reflect.Slice {
												var elem reflect.Value
												switch typ := rv.Field(j).Field(k).Type().Elem(); typ.Kind() {
												case reflect.Ptr:
													elem = reflect.New(typ.Elem())
												case reflect.Struct:
													elem = reflect.New(typ).Elem()
												}
												for l := 0; l < elem.NumField(); l++ {
													if !strings.EqualFold(elem.Type().Field(l).Name, slice_element[2]) {
														continue
													}
													if reflect.TypeOf(val) == elem.Field(l).Type() {
														elem.Field(l).Set(reflect.ValueOf(val))
													}
												}
												rv.Field(j).Field(k).Set(reflect.Append(rv.Field(j).Field(k), elem))
											}
											if reflect.TypeOf(val) == rv.Field(j).Field(k).Type() {
												rv.Field(j).Field(k).Set(reflect.ValueOf(val))
											}
										}
									}
								}
							}
						}
					}
					if len(slice_element) == 1 {
						if strings.EqualFold(rv.Type().Field(j).Name, slice_element[0]) {
							switch val.(type) {
							case int16, int32:
								if rv.Field(j).Kind() == reflect.Ptr {
									p := reflect.New(reflect.TypeOf(val))
									p.Elem().Set(reflect.ValueOf(val))
									rv.Field(j).Set(p)
								} else {
									rv.Field(j).Set(reflect.ValueOf(val))
								}
							case pgtype.TextArray:
								s := val.(pgtype.TextArray)
								var s_arr []string
								s.AssignTo(&s_arr)
								if reflect.TypeOf(s_arr) == rv.Field(j).Type() {
									rv.Field(j).Set(reflect.ValueOf(s_arr))
								}
							case pgtype.Numeric:
								n := val.(pgtype.Numeric)
								var nf float64
								n.AssignTo(&nf)
								switch rv.Field(j).Kind() {
								case reflect.Ptr:
									if reflect.TypeOf(nf) == rv.Field(j).Type() {
										rv.Field(j).Set(reflect.ValueOf(nf))
									} else if rv.Field(j).Type().Elem() == reflect.TypeOf(nf) {
										rv.Field(j).Set(reflect.ValueOf(&nf))
									}
								case reflect.Struct:
									for k := 0; k < rv.Field(j).NumField(); k++ {
										if reflect.TypeOf(nf) == rv.Field(j).Field(k).Type() {
											rv.Field(j).Field(k).Set(reflect.ValueOf(nf))
										}
									}
								default:
									if reflect.TypeOf(nf) == rv.Field(j).Type() {
										rv.Field(j).Set(reflect.ValueOf(nf))
									}
								}
							case time.Time:
								t, ok := val.(time.Time)
								if ok {
									rv.Field(j).Set(reflect.ValueOf(t))
								}
							default:
								if reflect.TypeOf(val) == rv.Field(j).Type() && strings.EqualFold(fieldName, slice_element[0]) {
									rv.Field(j).Set(reflect.ValueOf(val))
								}
							}
						}
					}
				}
			}
			return 1
		}
	}
	return 0
}
