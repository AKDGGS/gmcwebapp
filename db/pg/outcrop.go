package pg

import (
	dbf "gmc/db/flag"
)

func (pg *Postgres) GetOutcrop(id int, flags int) (map[string]interface{}, error) {
	outcrop, err := pg.queryRow("pg/outcrop_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if outcrop == nil {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file_byoutcropid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			outcrop["files"] = files
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_byoutcropid.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			outcrop["inventory"] = inventory
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/outcrop_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			outcrop["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_byoutcropid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			outcrop["quadrangles"] = qds
		}
	}

	return outcrop, nil
}
