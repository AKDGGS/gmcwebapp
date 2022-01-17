package db

func (pg *Postgres) GetBorehole(id int, flags int) (map[string]interface{}, error) {
	borehole, err := pg.queryRow("pg/borehole_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if borehole == nil {
		return nil, nil
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

	return borehole, nil
}
