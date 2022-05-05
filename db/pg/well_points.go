package pg

func (pg *Postgres) GetWellPoints() ([]map[string]interface{}, error) {
	return pg.queryRows("pg/well_points.sql")
}
