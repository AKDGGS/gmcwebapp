package pg

func (pg *Postgres) GetWellsPoints() (map[string]interface{}, error) {
	pl, err := pg.queryRows("pg/wells_points.sql")
	if err != nil {
		return nil, err
	}

	if pl == nil {
		return nil, nil
	}

	wellPoints := make(map[string]interface{})
	wellPoints["points"] = pl

	return wellPoints, nil
}
