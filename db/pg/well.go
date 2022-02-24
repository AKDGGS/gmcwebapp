package pg

import (
	dbf "gmc/db/flag"
)

func (pg *Postgres) GetWell(id int, flags int) (map[string]interface{}, error) {
	well, err := pg.queryRow("pg/well_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if well == nil {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			well["files"] = files
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_bywellid.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			well["inventory"] = inventory
		}
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		operators, err := pg.queryRows("pg/organization_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if operators != nil {
			well["operators"] = operators
		}
	}

	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			well["urls"] = urls
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/well_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			well["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			well["quadrangles"] = qds
		}
	}

	return well, nil
}
