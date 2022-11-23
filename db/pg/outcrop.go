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

	oc := model.Outcrop{}
	rowToStruct(rows, &oc)

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, _ := pg.pool.Query(context.Background(), q, id)
		rowToStruct(r, &oc.File)
	}
	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		kw, err := pg.queryRows(
			"pg/keyword/group_by_outcrop_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if kw != nil {
			for _, v := range kw {
				oc.Keywords = append(oc.Keywords, v)
			}
		}
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		q, err = assets.ReadString("pg/organization/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, _ := pg.pool.Query(context.Background(), q, id)
		rowToStruct(r, &oc.Organization)
	}

	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, _ := pg.pool.Query(context.Background(), q, id)
		rowToStruct(r, &oc.URL)
	}

	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, _ := pg.pool.Query(context.Background(), q, id)
		rowToStruct(r, &oc.Note)
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/outcrop/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		oc.GeoJSON = geojson["geojson"].(map[string]interface{})
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_outcrop_id.sql")
		if err != nil {
			return nil, err
		}
		r, _ := pg.pool.Query(context.Background(), q, id)
		rowToStruct(r, &oc.Quadrangle)
	}
	return &oc, nil
}
