package pg

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gmc/assets"
	"net/url"
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
