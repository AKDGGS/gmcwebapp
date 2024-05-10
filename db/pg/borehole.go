package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetBorehole(id int, flags int) (*model.Borehole, error) {
	q, err := assets.ReadString("pg/borehole/by_borehole_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	borehole := model.Borehole{}

	c, err := rowsToStruct(rows, &borehole)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.Files)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		q, err = assets.ReadString("pg/keyword/group_by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id, ((flags & dbf.PRIVATE) == 0))
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.KeywordSummary)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		q, err = assets.ReadString("pg/organization/by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.Organizations)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.URLs)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.Notes)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryValue("pg/borehole/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		borehole.GeoJSON = geojson
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		q, err = assets.ReadString("pg/mining_district/by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.MiningDistricts)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_borehole_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &borehole.Quadrangles)
		if err != nil {
			return nil, err
		}
	}
	return &borehole, nil
}
