package pg

import "fmt"

func (pg *Postgres) GetWellPoints() ([]map[string]interface{}, error) {
	fmt.Println(pg.queryRows("pg/well_points.sql"))
	return pg.queryRows("pg/well_points.sql")
}
