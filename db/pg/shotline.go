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

	if c, err := rowsToStruct(rows, &shotline); err != nil || c == 0 {
		return nil, err
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
		_, err = rowsToStruct(r, &shotline.Shotpoints)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(r, &shotline.KeywordSummary)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(r, &shotline.URLs)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(r, &shotline.Notes)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(r, &shotline.Quadrangles)
		if err != nil {
			return nil, err
		}
	}
	return &shotline, nil
}
