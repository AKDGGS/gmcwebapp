package pg

import dbf "gmc/db/flag"

func (pg *Postgres) GetProspect(id int, flags int) (map[string]interface{}, error) {
	prospect, err := pg.queryRow("pg/prospect/by_prospect_id.sql", id)
	if err != nil {
		return nil, err
	}
	if prospect == nil {
		return nil, nil
	}

	boreholes, err := pg.queryRows("pg/borehole/by_prospect_id.sql", id)
	if err != nil {
		return nil, err
	}
	if boreholes != nil {
		prospect["boreholes"] = boreholes
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file/by_prospect_id.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			prospect["files"] = files
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		kw, err := pg.queryRows(
			"pg/keyword/group_by_prospect_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if kw != nil {
			prospect["keywords"] = kw
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/prospect/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			prospect["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		mds, err := pg.queryRows("pg/mining_district/by_prospect_id.sql", id)
		if err != nil {
			return nil, err
		}
		if mds != nil {
			prospect["mining_districts"] = mds
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle/250k_by_prospect_id.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			prospect["quadrangles"] = qds
		}
	}

	return prospect, nil
}
