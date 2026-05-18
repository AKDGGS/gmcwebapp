package pg

func (pg *Postgres) Clean() error {
	err := pg.executeSQL("pg/inventory/delete_inactive.sql")
	if err != nil {
		return err
	}
	return nil
}
