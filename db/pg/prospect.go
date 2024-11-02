package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) ListProspects() ([]*model.Prospect, error) {
	q, err := assets.ReadString("pg/prospect/list.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	prospects, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.Prospect])
	if err != nil {
		return nil, err
	}
	return prospects, nil
}

func (pg *Postgres) GetProspect(id int, flags int) (*model.Prospect, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	prospect, err := cQryStruct[model.Prospect](
		conn, "pg/prospect/by_prospect_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no prospect is found, stop right here
	if prospect == nil {
		return nil, nil
	}

	if (flags & dbf.BOREHOLE) != 0 {
		prospect.Boreholes, err = cQryStructs[model.Borehole](
			conn, "pg/borehole/by_prospect_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.FILES) != 0 {
		prospect.Files, err = cQryStructs[model.File](
			conn, "pg/file/by_prospect_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		prospect.KeywordSummary, err = cQryStructs[model.KeywordSummary](
			conn, "pg/keyword/group_by_prospect_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		prospect.MiningDistricts, err = cQryStructs[model.MiningDistrict](
			conn, "pg/mining_district/by_prospect_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		prospect.Quadrangles, err = cQryStructs[model.Quadrangle](
			conn, "pg/quadrangle/250k_by_prospect_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		prospect.GeoJSON, err = cQryValue[interface{}](
			conn, "pg/prospect/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}
	return prospect, nil
}
