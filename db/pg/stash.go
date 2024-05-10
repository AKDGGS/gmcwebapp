package pg

func (pg *Postgres) GetStash(id int) (interface{}, error) {
	return pg.queryValue("pg/inventory/stash.sql", id)
}
