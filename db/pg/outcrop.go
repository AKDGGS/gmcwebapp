package pg

import (
	"context"

	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetOutcrop(id int, flags int) (*model.Outcrop, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	outcrop, err := cQryStruct[model.Outcrop](
		conn, "pg/outcrop/by_outcrop_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no outcrop is found, stop right here
	if outcrop == nil {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		outcrop.Files, err = cQryStructs[model.File](
			conn, "pg/file/by_outcrop_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		outcrop.KeywordSummary, err = cQryStructs[model.KeywordSummary](
			conn, "pg/keyword/group_by_outcrop_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		outcrop.URLs, err = cQryStructs[model.URL](
			conn, "pg/url/by_outcrop_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		outcrop.Notes, err = cQryStructs[model.Note](
			conn, "pg/note/by_outcrop_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		outcrop.Quadrangles, err = cQryStructs[model.Quadrangle](
			conn, "pg/quadrangle/250k_by_outcrop_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		outcrop.GeoJSON, err = cQryValue(
			conn, "pg/outcrop/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	return outcrop, nil
}
