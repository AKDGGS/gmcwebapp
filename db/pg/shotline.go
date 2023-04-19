package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetShotline(id int, flags int) (*model.Shotline, error) {
	q, err := assets.ReadString("pg/shotline/by_shotline_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shotline := model.Shotline{}

	if rowToStruct(rows, &shotline) == 0 {
		return nil, nil
	}
	if (flags & dbf.SHOTPOINT) != 0 {
		q, err = assets.ReadString("pg/shotpoint/by_shotline_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &shotline.Shotpoints)
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		q, err = assets.ReadString("pg/keyword/group_by_shotline_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id, ((flags & dbf.PRIVATE) == 0))
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &shotline.KeywordSummary)
	}
	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_shotline_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &shotline.URLs)
	}
	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_shotline_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &shotline.Notes)
	}
	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/shotline/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson["geojson"] != nil {
			shotline.GeoJSON = geojson["geojson"]
		}
	}
	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_shotline_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		rowToStruct(r, &shotline.Quadrangles)
	}
	return &shotline, nil
}
