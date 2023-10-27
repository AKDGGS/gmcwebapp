package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetProspect(id int, flags int) (*model.Prospect, error) {
	q, err := assets.ReadString("pg/prospect/by_prospect_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	prospect := model.Prospect{}

	if count, err := rowsToStruct(rows, &prospect); err != nil || count == 0 {
		return nil, err
	}
	if (flags & dbf.BOREHOLE) != 0 {
		q, err := assets.ReadString("pg/borehole/by_prospect_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &prospect.Boreholes)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_prospect_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &prospect.Files)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		q, err = assets.ReadString("pg/keyword/group_by_prospect_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id, ((flags & dbf.PRIVATE) == 0))
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &prospect.KeywordSummary)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/prospect/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson["geojson"] != nil {
			prospect.GeoJSON = geojson["geojson"]
		}
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		q, err = assets.ReadString("pg/mining_district/by_prospect_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &prospect.MiningDistricts)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		q, err = assets.ReadString("pg/quadrangle/250k_by_prospect_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		_, err = rowsToStruct(r, &prospect.Quadrangles)
		if err != nil {
			return nil, err
		}
	}
	return &prospect, nil
}
