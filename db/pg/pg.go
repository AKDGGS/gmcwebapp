package pg

import (
	"context"
	"fmt"
	"net/url"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRows error = pgx.ErrNoRows

type Postgres struct {
	pool *pgxpool.Pool
}

func New(cfg map[string]interface{}) (*Postgres, error) {
	u, ok := cfg["url"].(string)
	if !ok {
		return nil, fmt.Errorf("database url required")
	}

	dburl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if dburl.Scheme != "postgres" && dburl.Scheme != "postgresql" {
		return nil, fmt.Errorf(
			"postgres database urls must have scheme " +
				"of \"postgres\" or \"postgresql\"",
		)
	}

	pgcfg, err := pgxpool.ParseConfig(dburl.String())
	if err != nil {
		return nil, err
	}

	if con_min, ok := cfg["min_connections"].(int); ok {
		pgcfg.MinConns = int32(con_min)
	}

	if con_max, ok := cfg["max_connections"].(int); ok {
		if int32(con_max) < pgcfg.MinConns {
			return nil, fmt.Errorf(
				"max_connections must be greater than or equal "+
					"to min_connections (%d <= %d)", pgcfg.MinConns, con_max,
			)
		}
		pgcfg.MaxConns = int32(con_max)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgcfg)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (pg *Postgres) Shutdown() {
	pg.pool.Close()
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

// Using the provided connection, runs the query found at the given
// asset location, and returns the results via pgx.CollectRows() as []T
func cQryStructs[T any](conn *pgxpool.Conn, qry string, args ...interface{}) ([]T, error) {
	sql, err := assets.ReadString(qry)
	if err != nil {
		return nil, err
	}
	r, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", qry, err)
	}
	v, err := pgx.CollectRows(r, pgx.RowToStructByNameLax[T])
	if err != nil {
		return nil, fmt.Errorf("%s: %s", qry, err)
	}
	return v, nil
}

// Using the provided connection, runs the query found at the given
// asset location, and returns the results via pgx.CollectOneRow() as []T
func cQryStruct[T any](conn *pgxpool.Conn, qry string, args ...interface{}) (*T, error) {
	sql, err := assets.ReadString(qry)
	if err != nil {
		return nil, err
	}
	r, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", qry, err)
	}
	v, err := pgx.CollectOneRow(r, pgx.RowToAddrOfStructByNameLax[T])
	if err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %s", qry, err)
	}
	return v, nil
}

// Using the provided connection, runs the query found at the given
// asset location, and then returns the singular result
func cQryValue(conn *pgxpool.Conn, qry string, args ...interface{}) (interface{}, error) {
	sql, err := assets.ReadString(qry)
	if err != nil {
		return nil, err
	}

	var x interface{}
	row := conn.QueryRow(context.Background(), sql, args...)
	if err := row.Scan(&x); err != nil && err != ErrNoRows {
		return nil, fmt.Errorf("%s: %s", qry, err)
	}
	return x, nil
}
