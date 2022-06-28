package pg

func (pg *Postgres) GetStash(id int) (map[string]interface{}, error) {
	return pg.queryRow("pg/inventory/stash.sql", id)
}
