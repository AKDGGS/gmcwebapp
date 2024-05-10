package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetWell(id int, flags int) (*model.Well, error) {
	q, err := assets.ReadString("pg/well/by_well_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	well := model.Well{}

	c, err := rowsToStruct(rows, &well)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.Files)
		if err != nil {
			return nil, err
		}
	}
	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		q, err = assets.ReadString("pg/keyword/group_by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id, ((flags & dbf.PRIVATE) == 0))
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.KeywordSummary)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		q, err = assets.ReadString("pg/organization/by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.Organizations)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.URLs)
		if err != nil {
			return nil, err
		}
	}
	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.Notes)
		if err != nil {
			return nil, err
		}
	}
	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryValue("pg/well/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		well.GeoJSON = geojson
	}
	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_well_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &well.Quadrangles)
		if err != nil {
			return nil, err
		}
	}
	return &well, nil
}
