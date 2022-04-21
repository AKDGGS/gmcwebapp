package pg

func (pg *Postgres) GetWellsPointList() (map[string]interface{}, error) {
	pl, err := pg.queryRows("pg/well_pointlist.sql")
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
