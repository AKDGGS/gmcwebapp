package db

func (pg *Postgres) GetBorehole(id int, flags int) (map[string]interface{}, error) {
	borehole, err := pg.queryRow("pg/borehole_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if borehole == nil {
		return nil, nil
	}

	if (flags & FILES) != 0 {
		files, err := pg.queryRows("pg/file_byboreholeid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			borehole["files"] = files
		}
	}

	if (flags & INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_byboreholeid.sql", id,
			((flags & PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			borehole["inventory"] = inventory
		}
	}

	if (flags & GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/borehole_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			borehole["geojson"] = geojson["geojson"]
		}
	}

	if (flags & MINING_DISTRICTS) != 0 {
		mds, err := pg.queryRows("pg/miningdistrict_byboreholeid.sql", id)
		if err != nil {
			return nil, err
		}
		if mds != nil {
			borehole["mining_districts"] = mds
		}
	}

	if (flags & QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_byboreholeid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			borehole["quadrangles"] = qds
		}
	}

	return borehole, nil
}
