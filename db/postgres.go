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
	/*
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
	*/
	return 0, "", time.Time{}, nil
}

func (pg *Postgres) GetProspect(id int, flags int) (map[string]interface{}, error) {
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

	if (flags & FILES) != 0 {
		files, err := pg.queryRows("pg/file_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			prospect["files"] = files
		}
	}

	if (flags & INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_byprospectid.sql", id,
			((flags & SHOW_PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			prospect["inventory"] = inventory
		}
	}

	if (flags & GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/prospect_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			prospect["geojson"] = geojson["geojson"]
		}
	}

	if (flags & MINING_DISTRICTS) != 0 {
		mds, err := pg.queryRows("pg/miningdistrict_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if mds != nil {
			prospect["mining_districts"] = mds
		}
	}

	if (flags & QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			prospect["quadrangles"] = qds
		}
	}

	return prospect, nil
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
