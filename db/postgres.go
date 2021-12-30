package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gmc/assets"
	"net/url"
	"time"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func newPostgres(u *url.URL) (*Postgres, error) {
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

func (pg *Postgres) GetFile(id int, include_private bool) (int, string, time.Time, error) {
	var q string
	if include_private {
		q = assets.ReadString("pg/file_all.sql")
	} else {
		q = assets.ReadString("pg/file_public.sql")
	}

	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return 0, "", time.Time{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var aid int
		var fname string
		var ftime time.Time
		err = rows.Scan(&aid, &fname, &ftime)
		if err != nil {
			return 0, "", time.Time{}, err
		}
		return aid, fname, ftime, nil
	} else {
		return 0, "", time.Time{}, nil
	}
}

func (pg *Postgres) GetProspect(id int) (map[string]interface{}, error) {
	prospect, err := pg.queryRow("pg/prospect_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if prospect == nil {
		return nil, nil
	}

	boreholes, err := pg.queryRows("pg/borehole_byprospectid.sql", id)
	if err != nil {
		return nil, err
	}
	if boreholes != nil {
		prospect["boreholes"] = boreholes
	}

	files, err := pg.queryRows("pg/file_byprospectid.sql", id)
	if err != nil {
		return nil, err
	}
	if files != nil {
		prospect["files"] = files
	}

	geojson, err := pg.queryRow("pg/prospect_geojson.sql", id)
	if err != nil {
		return nil, err
	}
	if geojson != nil {
		prospect["geojson"] = geojson["geojson"]
	}

	mds, err := pg.queryRows("pg/miningdistrict_byprospectid.sql", id)
	if err != nil {
		return nil, err
	}
	if mds != nil {
		prospect["mining_districts"] = mds
	}

	qds, err := pg.queryRows("pg/quadrangle250k_byprospectid.sql", id)
	if err != nil {
		return nil, err
	}
	if qds != nil {
		prospect["quadrangles"] = qds
	}

	return prospect, nil
}

func (pg *Postgres) queryRow(name string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := pg.pool.Query(
		context.Background(), assets.ReadString(name), args...,
	)
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
	rows, err := pg.pool.Query(
		context.Background(), assets.ReadString(name), args...,
	)
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
