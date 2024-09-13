package pg

import (
	"context"

	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetShotline(id int, flags int) (*model.Shotline, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	shotline, err := cQryStruct[model.Shotline](
		conn, "pg/shotline/by_shotline_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no shotline is found, stop right here
	if shotline == nil {
		return nil, nil
	}

	if (flags & dbf.SHOTPOINT) != 0 {
		shotline.Shotpoints, err = cQryStructs[model.Shotpoint](
			conn, "pg/shotpoint/by_shotline_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		shotline.KeywordSummary, err = cQryStructs[model.KeywordSummary](
			conn, "pg/keyword/group_by_shotline_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		shotline.URLs, err = cQryStructs[model.URL](
			conn, "pg/url/by_shotline_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		shotline.Notes, err = cQryStructs[model.Note](
			conn, "pg/note/by_shotline_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		shotline.Quadrangles, err = cQryStructs[model.Quadrangle](
			conn, "pg/quadrangle/250k_by_shotline_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		shotline.GeoJSON, err = cQryValue[interface{}](
			conn, "pg/shotline/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}
	return shotline, nil
}
