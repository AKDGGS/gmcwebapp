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
			for i, val := range columnValues {
				fieldName := string(r.FieldDescriptions()[i].Name)
				for j := 0; j < rv.NumField(); j++ {

					if rv.Field(j).Kind() == reflect.Struct {
						// fmt.Println(rv.Field(j).Type().Name(), rv.Field(j).Kind())
						rowToStruct(r, rv.Field(j).Addr().Interface())
					}
					if !strings.EqualFold(fieldName, rv.Type().Field(j).Name) {
						continue
					}
					fmt.Println(rv.Type().Field(j).Name, rv.Type().Field(j).Name)
					switch val.(type) {
					case pgtype.TextArray:
						s := val.(pgtype.TextArray)
						var s_arr []string
						s.AssignTo(&s_arr)
						if reflect.TypeOf(s_arr) == rv.Field(j).Type() {
							rv.Field(j).Set(reflect.ValueOf(s_arr))
						}
					case pgtype.Numeric:
						n := val.(pgtype.Numeric)
						var nv float64
						n.AssignTo(&nv)
						if reflect.TypeOf(nv) == rv.Field(j).Type() {
							rv.Field(j).Set(reflect.ValueOf(nv))
						}
					case time.Time:
						t, ok := val.(time.Time)
						if ok {
							rv.Field(j).Set(reflect.ValueOf(&t))
						}
					default:
						if reflect.TypeOf(val) == rv.Field(j).Type() {
							rv.Field(j).Set(reflect.ValueOf(val))
						}
					}
				}
			}
			return 1
		}
	}
	return 0
}

func ReadStruct(r pgx.Rows, st interface{}) {
	var rowResults = make(map[string]interface{})
	for r.Next() {
		columnValues, _ := r.Values()
		for i, val := range columnValues {
			fieldName := string(r.FieldDescriptions()[i].Name)
			rowResults[fieldName] = val
		}
	}
	readStruct(r, reflect.ValueOf(st), rowResults)
}

func readStruct(r pgx.Rows, rv reflect.Value, rowResults map[string]interface{}) int {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.CanSet() {
		return 0
	}
	fmt.Println("==============" + rv.Type().Name() + "===============")

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
			if c := readStruct(r, elem.Addr(), rowResults); c > 0 {
				rowCount += c
				rv.Set(reflect.Append(rv, elem))
				continue
			}
			return rowCount
		}
	case reflect.Struct:
		for j := 0; j < rv.NumField(); j++ {
			switch rv.Field(j).Kind() {
			case reflect.Struct:
				readStruct(r, rv.Field(j), rowResults)
			default:
				if reflect.TypeOf(rowResults[rv.Type().Field(j).Name]) == rv.Field(j).Type() {
					fmt.Println(rv.Type().Field(j).Name, rowResults[rv.Type().Field(j).Name])
					rv.Field(j).Set(reflect.ValueOf(rowResults[rv.Type().Field(j).Name]))

					// Line 304 isn't setting the value.  I think I need to pass the value as a pointer or Addr().Interface()
				}
			}
		}
	}
	return 0
}

//Some of this might be useful
// func ReadStruct(r pgx.Rows, st interface{}) {
// 	readStruct(r, reflect.ValueOf(st))
// }
//
// func readStruct(r pgx.Rows, rv reflect.Value) int {
// 	if rv.Kind() == reflect.Ptr {
// 		rv = rv.Elem()
// 	}
// 	if !rv.CanSet() {
// 		return 0
// 	}
// 	fmt.Println("==============" + rv.Type().Name() + "===============")
// 	switch rv.Kind() {
// 	case reflect.Slice:
// 		var elem reflect.Value
// 		switch typ := rv.Type().Elem(); typ.Kind() {
// 		case reflect.Ptr:
// 			elem = reflect.New(typ.Elem())
// 		case reflect.Struct:
// 			elem = reflect.New(typ).Elem()
// 		}
// 		ReadStruct(r, elem.Addr().Interface())
// 		for rowCount := 0; ; {
// 			ReadStruct(r, elem.Addr().Interface())
// 			rv.Set(reflect.Append(rv, elem))
// 			// if c := rowToStruct(r, elem.Addr().Interface()); c > 0 {
// 			// 	rowCount += c
// 			// 	rv.Set(reflect.Append(rv, elem))
// 			// 	continue
// 			// }
// 			return rowCount
// 		}
// 	case reflect.Struct:
// 		fmt.Println("============= Struct", rv.Type().Name(), rv.NumField(), "=================")
// 		for j := 0; j < rv.NumField(); j++ {
// 			if r.Next() {
// 				columnValues, _ := r.Values()
// 				for i, _ := range columnValues {
// 					fieldName := string(r.FieldDescriptions()[i].Name)
// 					// fmt.Println(fieldName, val)
// 					if strings.Contains(fieldName, ".") {
// 						sArr := strings.Split(fieldName, ".")
// 						fmt.Println(rv.Type().Name(), sArr[0], sArr[1])
// 						if (rv.FieldByName(sArr[0]) != reflect.Value{}) {
// 							fmt.Println(reflect.ValueOf(rv.FieldByName(sArr[0]).Type()))
// 							ReadStruct(r, rv.FieldByName(sArr[0]))
// 						}
// 					}
// 					// else {
// 					// 	fmt.Println(rv.Field(j).Type().Name(), rv.Type().Name(), val)
// 					// }
// 				}
//
// 				// for i, val := range columnValues {
// 				// 	fieldName := string(r.FieldDescriptions()[i].Name)
// 				// 	fmt.Println(fieldName, val)
// 				// 	fmt.Println("-------- Fields: ", rv.Field(j).Type().Name())
// 				// }
// 			}
// 		}
// 	}
//
// 	// for j := 0; j < rv.NumField(); j++ {
// 	// 	// fmt.Println("Fields: "+rv.Field(j).Type().Name(), rv.Field(j).Kind())
// 	// 	if rv.Field(j).Kind() == reflect.Struct {
// 	// 		fmt.Println("====", rv.Type().Name())
// 	// 		typ := rv.Type().Elem()
// 	// 		elem := reflect.New(typ).Elem()
// 	// 		ReadStruct(r, elem.Addr().Interface())
// 	// 		rv.Field(j).Set(elem)
// 	// 	}
// 	// else {
// 	// 	if r.Next() {
// 	// 		columnValues, _ := r.Values()
// 	// 		for i, val := range columnValues {
// 	// 			fieldName := string(r.FieldDescriptions()[i].Name)
// 	// 			for j := 0; j < rv.NumField(); j++ {
// 	// 				f := rv.Field(j)
// 	// 				switch f.Kind() {
// 	// 				case reflect.Struct:
// 	// 					fmt.Println(">>>>>>>>>>>>>>", f.Kind(), "|", f.Type(), "|", f.Type().Name())
// 	// 					readStruct(r, f)
// 	// 				case reflect.Slice:
// 	// 					for k := 0; k < f.Len(); k++ {
// 	// 						readStruct(r, f.Index(k))
// 	// 					}
// 	// 				default:
// 	// 					// fmt.Println(fieldName, rv.Type().Field(j).Name)
// 	// 					if !strings.EqualFold(fieldName, rv.Type().Field(j).Name) {
// 	// 						continue
// 	// 					}
// 	// 					// fmt.Println(fieldName, rv.Type().Field(j).Name, val)
// 	// 					switch val.(type) {
// 	// 					case pgtype.TextArray:
// 	// 						s := val.(pgtype.TextArray)
// 	// 						var s_arr []string
// 	// 						s.AssignTo(&s_arr)
// 	// 						if reflect.TypeOf(s_arr) == rv.Field(j).Type() {
// 	// 							rv.Field(j).Set(reflect.ValueOf(s_arr))
// 	// 						}
// 	// 					case pgtype.Numeric:
// 	// 						n := val.(pgtype.Numeric)
// 	// 						var nv float64
// 	// 						n.AssignTo(&nv)
// 	// 						if reflect.TypeOf(nv) == rv.Field(j).Type() {
// 	// 							rv.Field(j).Set(reflect.ValueOf(nv))
// 	// 						}
// 	// 					case time.Time:
// 	// 						t, ok := val.(time.Time)
// 	// 						if ok {
// 	// 							rv.Field(j).Set(reflect.ValueOf(&t))
// 	// 						}
// 	// 					default:
// 	// 						if reflect.TypeOf(val) == rv.Field(j).Type() {
// 	// 							// fmt.Println(rv.Type().Field(j).Name, val)
// 	// 							rv.Field(j).Set(reflect.ValueOf(val))
// 	// 						}
// 	// 					}
// 	// 				}
// 	// 				// case reflect.String:
// 	// 				// 	fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 				// 	val.Field(i).SetString("Test")
// 	// 				// case reflect.Int32:
// 	// 				// 	fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 				// 	val.Field(i).SetInt(99)
// 	// 				// 	// fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 				// default:
// 	// 				// 	if !strings.EqualFold(fieldName, val.Type().Field(i).Name) {
// 	// 				// 		continue
// 	// 				// 	}
// 	// 				// 	if reflect.TypeOf(val) == rv.Field(j).Type() {
// 	// 				// 		// fmt.Println(rv.Type().Field(j).Name, val)
// 	// 				// 		rv.Field(j).Set(reflect.ValueOf(val))
// 	// 				// 	}
// 	// 			}
// 	// 		}
// 	// 		return 1
// 	// 	}
// 	// }
// 	// }
//
// 	// if r.Next() {
// 	// 	columnValues, _ := r.Values()
// 	// 	for i, val := range columnValues {
// 	// 		fieldName := string(r.FieldDescriptions()[i].Name)
// 	// 		// fmt.Println(fieldName, val)
// 	// 		for j := 0; j < rv.NumField(); j++ {
// 	// 			f := rv.Field(j)
// 	// 			switch f.Kind() {
// 	// 			case reflect.Struct:
// 	// 				fmt.Println(">>>>>>>>>>>>>>", f.Kind(), "|", f.Type(), "|", f.Type().Name())
// 	// 				readStruct(r, f)
// 	// 			case reflect.Slice:
// 	// 				for k := 0; k < f.Len(); k++ {
// 	// 					readStruct(r, f.Index(k))
// 	// 				}
// 	// 			default:
// 	// 				// fmt.Println(fieldName, rv.Type().Field(j).Name)
// 	// 				if !strings.EqualFold(fieldName, rv.Type().Field(j).Name) {
// 	// 					continue
// 	// 				}
// 	// 				// fmt.Println(fieldName, rv.Type().Field(j).Name, val)
// 	// 				switch val.(type) {
// 	// 				case pgtype.TextArray:
// 	// 					s := val.(pgtype.TextArray)
// 	// 					var s_arr []string
// 	// 					s.AssignTo(&s_arr)
// 	// 					if reflect.TypeOf(s_arr) == rv.Field(j).Type() {
// 	// 						rv.Field(j).Set(reflect.ValueOf(s_arr))
// 	// 					}
// 	// 				case pgtype.Numeric:
// 	// 					n := val.(pgtype.Numeric)
// 	// 					var nv float64
// 	// 					n.AssignTo(&nv)
// 	// 					if reflect.TypeOf(nv) == rv.Field(j).Type() {
// 	// 						rv.Field(j).Set(reflect.ValueOf(nv))
// 	// 					}
// 	// 				case time.Time:
// 	// 					t, ok := val.(time.Time)
// 	// 					if ok {
// 	// 						rv.Field(j).Set(reflect.ValueOf(&t))
// 	// 					}
// 	// 				default:
// 	// 					if reflect.TypeOf(val) == rv.Field(j).Type() {
// 	// 						// fmt.Println(rv.Type().Field(j).Name, val)
// 	// 						rv.Field(j).Set(reflect.ValueOf(val))
// 	// 					}
// 	// 				}
// 	// 			}
// 	// 			// case reflect.String:
// 	// 			// 	fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 			// 	val.Field(i).SetString("Test")
// 	// 			// case reflect.Int32:
// 	// 			// 	fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 			// 	val.Field(i).SetInt(99)
// 	// 			// 	// fmt.Printf("%v=%v\n", val.Type().Field(i).Name, val.Field(i))
// 	// 			// default:
// 	// 			// 	if !strings.EqualFold(fieldName, val.Type().Field(i).Name) {
// 	// 			// 		continue
// 	// 			// 	}
// 	// 			// 	if reflect.TypeOf(val) == rv.Field(j).Type() {
// 	// 			// 		// fmt.Println(rv.Type().Field(j).Name, val)
// 	// 			// 		rv.Field(j).Set(reflect.ValueOf(val))
// 	// 			// 	}
// 	// 		}
// 	// 	}
// 	// 	return 1
// 	// }
//
// 	return 0
// }
