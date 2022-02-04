package pg

import dbf "gmc/db/flag"

func (pg *Postgres) GetProspect(id int, flags int) (map[string]interface{}, error) {
	prospect, err := pg.queryRow("pg/prospect_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if prospect == nil {
		return nil, nil
	}

	boreholes, err := pg.queryRows("pg/borehole_byprospectid.sql", id)
	if err != nil {
		return nil, err
	}
	if boreholes != nil {
		prospect["boreholes"] = boreholes
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			prospect["files"] = files
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_byprospectid.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			prospect["inventory"] = inventory
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/prospect_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			prospect["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.MINING_DISTRICTS) != 0 {
		mds, err := pg.queryRows("pg/miningdistrict_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if mds != nil {
			prospect["mining_districts"] = mds
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_byprospectid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			prospect["quadrangles"] = qds
		}
	}

	return prospect, nil
}
