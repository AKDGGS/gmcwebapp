package pg

func (pg *Postgres) GetWellPoints() ([]map[string]interface{}, error) {
	pts, err := pg.queryRows("pg/well_points.sql")
	if err != nil {
		return nil, err
	}

	if pts == nil {
		return nil, nil
	}

	return pts, nil
}
