package pg

import (
	"context"
	"fmt"
	"net/url"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

	config, err := pgxpool.ParseConfig(dburl.String())
	if err != nil {
		return nil, err
	}

	config.LazyConnect = true
	if lazyconnect, ok := cfg["lazyconnect"].(bool); ok {
		config.LazyConnect = lazyconnect
	}

	if con_min, ok := cfg["min_connections"].(int); ok {
		config.MinConns = int32(con_min)
	}

	if con_max, ok := cfg["max_connections"].(int); ok {
		if int32(con_max) < config.MinConns {
			return nil, fmt.Errorf(
				"max_connections must be greater than or equal "+
					"to min_connections (%d <= %d)", config.MinConns, con_max,
			)
		}
		config.MaxConns = int32(con_max)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (pg *Postgres) Shutdown() {
	pg.pool.Close()
}

func (pg *Postgres) queryValue(name string, args ...interface{}) (interface{}, error) {
	q, err := assets.ReadString(name)
	if err != nil {
		return nil, err
	}

	var val interface{}
	row := pg.pool.QueryRow(context.Background(), q, args...)
	if err := row.Scan(&val); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
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

// Using the provided connection, runs the query found at the given asset
// location and returns the results into dest via rowsToStruct()
func ConnQuery(conn *pgxpool.Conn, qry string, dest interface{}, args ...interface{}) (int, error) {
	sql, err := assets.ReadString(qry)
	if err != nil {
		return 0, err
	}

	r, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		return 0, err
	}

	c, err := rowsToStruct(r, dest)
	r.Close()
	return c, err
}
