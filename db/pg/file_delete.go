package pg

import (
	"context"
)

func (pg *Postgres) DeleteFile(id int, rm_links bool) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if rm_links {
		err = pg.executeSQL("pg/file/delete_borehole_link_by_file_id.sql", id)
		if err != nil {
			return err
		}
		err = pg.executeSQL("pg/file/delete_well_link_by_file_id.sql", id)
		if err != nil {
			return err
		}
		err = pg.executeSQL("pg/file/delete_inventory_link_by_file_id.sql", id)
		if err != nil {
			return err
		}
		err = pg.executeSQL("pg/file/delete_outcrop_link_by_file_id.sql", id)
		if err != nil {
			return err
		}
		err = pg.executeSQL("pg/file/delete_prospect_link_by_file_id.sql", id)
		if err != nil {
			return err
		}
	}

	err = pg.executeSQL("pg/file/delete_by_file_id.sql", id)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
