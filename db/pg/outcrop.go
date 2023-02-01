package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetOutcrop(id int, flags int) (*model.Outcrop, error) {
	q, err := assets.ReadString("pg/outcrop/by_outcrop_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	outcrop := model.Outcrop{}
	rowToStruct(rows, &outcrop)
	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.Files)
	}
	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		q, err = assets.ReadString("pg/keyword/group_by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id, ((flags & dbf.PRIVATE) == 0))
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.KeywordSummary)
	}
	if (flags & dbf.ORGANIZATION) != 0 {
		q, err = assets.ReadString("pg/organization/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.Organizations)
	}
	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.URLs)
	}

	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.Notes)
	}
	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/outcrop/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson["geojson"] != nil {
			outcrop.GeoJSON = geojson["geojson"].(map[string]interface{})
		}
	}
	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &outcrop.Quadrangles)
	}
	return &outcrop, nil
}
