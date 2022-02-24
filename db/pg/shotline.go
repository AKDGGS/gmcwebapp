package pg

import dbf "gmc/db/flag"

func (pg *Postgres) GetShotline(id int, flags int) (map[string]interface{}, error) {
	shotline, err := pg.queryRow("pg/shotline_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if shotline == nil {
		return nil, nil
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_byshotlineid.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			shotline["inventory"] = inventory
		}
	}
	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url_byshotlineid.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			shotline["urls"] = urls
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/shotline_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			shotline["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_byshotlineid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			shotline["quadrangles"] = qds
		}
	}

	return shotline, nil
}
