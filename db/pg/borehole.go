package pg

import (
	"context"

	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetBorehole(id int, flags int) (*model.Borehole, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	borehole, err := cQryStruct[model.Borehole](
		conn, "pg/borehole/by_borehole_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no borehole is found, stop right here
	if borehole == nil {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		borehole.Files, err = cQryStructs[model.File](
			conn, "pg/file/by_borehole_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		borehole.KeywordSummary, err = cQryStructs[model.KeywordSummary](
			conn, "pg/keyword/group_by_borehole_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		borehole.URLs, err = cQryStructs[model.URL](
			conn, "pg/url/by_borehole_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		borehole.Notes, err = cQryStructs[model.Note](
			conn, "pg/note/by_borehole_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		borehole.MiningDistricts, err = cQryStructs[model.MiningDistrict](
			conn, "pg/mining_district/by_borehole_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		borehole.Quadrangles, err = cQryStructs[model.Quadrangle](
			conn, "pg/quadrangle/250k_by_borehole_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		borehole.GeoJSON, err = cQryValue(
			conn, "pg/borehole/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	return borehole, nil
}
