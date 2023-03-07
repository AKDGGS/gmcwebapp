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
	rowToStruct(rows, &prospect)

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
		rowToStruct(rows, &prospect.Boreholes)
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
		rowToStruct(r, &prospect.Files)
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
		rowToStruct(r, &prospect.KeywordSummary)
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
		rowToStruct(r, &prospect.MiningDistricts)
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
		rowToStruct(r, &prospect.Quadrangles)
	}

	return &prospect, nil
}
