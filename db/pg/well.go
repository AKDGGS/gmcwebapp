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

	if rowsToStruct(rows, &well) == 0 {
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
		rowsToStruct(r, &well.Files)
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
		rowsToStruct(r, &well.KeywordSummary)
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
		rowsToStruct(r, &well.Organizations)
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
		rowsToStruct(r, &well.URLs)
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
		rowsToStruct(r, &well.Notes)
	}
	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/well/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson["geojson"] != nil {
			well.GeoJSON = geojson["geojson"]
		}
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
		rowsToStruct(r, &well.Quadrangles)
	}
	return &well, nil
}
