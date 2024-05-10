package pg

func (pg *Postgres) GetInventoryStash(id int) (interface{}, error) {
	return pg.queryValue("pg/inventory/stash.sql", id)
}
