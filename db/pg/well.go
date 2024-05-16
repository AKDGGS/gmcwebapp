package pg

import (
	"context"

	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetWell(id int, flags int) (*model.Well, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	well, err := cQryStruct[model.Well](
		conn, "pg/well/by_well_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no well is found, stop right here
	if well == nil {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		well.Files, err = cQryStructs[model.File](
			conn, "pg/file/by_well_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		well.KeywordSummary, err = cQryStructs[model.KeywordSummary](
			conn, "pg/keyword/group_by_well_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		well.URLs, err = cQryStructs[model.URL](
			conn, "pg/url/by_well_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		well.Notes, err = cQryStructs[model.Note](
			conn, "pg/note/by_well_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		well.Quadrangles, err = cQryStructs[model.Quadrangle](
			conn, "pg/quadrangle/250k_by_well_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		well.GeoJSON, err = cQryValue(
			conn, "pg/well/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	return well, nil
}
