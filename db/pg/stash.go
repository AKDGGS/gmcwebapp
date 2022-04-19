package pg

func (pg *Postgres) GetStash(id int) (map[string]interface{}, error) {
	stash, err := pg.queryRow("pg/stash_byinventoryid.sql", id)
	if err != nil {
		return nil, err
	}

	if stash == nil {
		return nil, nil
	}

	return stash, nil
}
