package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) DeleteFile(file *model.File, rm_links bool) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	if rm_links {
		if len(file.BoreholeIDs) > 0 {
			rm_well_link_sql, err := assets.ReadString("pg/file/delete_borehole_link_by_file_id.sql")
			if err != nil {
				return err
			}
			fct, err := tx.Exec(context.Background(), rm_well_link_sql, file.ID)
			if err != nil {
				return err
			}
			if fct.RowsAffected() < 1 {
				return fmt.Errorf("Borehole link not found for file id %d", file.ID)
			}
		}
		if len(file.InventoryIDs) > 0 {
			rm_well_link_sql, err := assets.ReadString("pg/file/delete_inventory_link_by_file_id.sql")
			if err != nil {
				return err
			}

			fct, err := tx.Exec(context.Background(), rm_well_link_sql, file.ID)
			if err != nil {
				return err
			}
			if fct.RowsAffected() < 1 {
				return fmt.Errorf("Inventory link not found for file id %d", file.ID)
			}
		}
		if len(file.OutcropIDs) > 0 {
			rm_well_link_sql, err := assets.ReadString("pg/file/delete_outcrop_link_by_file_id.sql")
			if err != nil {
				return err
			}
			fct, err := tx.Exec(context.Background(), rm_well_link_sql, file.ID)
			if err != nil {
				return err
			}
			if fct.RowsAffected() < 1 {
				return fmt.Errorf("Outcrop link not found for file id %d", file.ID)
			}
		}
		if len(file.ProspectIDs) > 0 {
			rm_well_link_sql, err := assets.ReadString("pg/file/delete_prospect_link_by_file_id.sql")
			if err != nil {
				return err
			}
			fct, err := tx.Exec(context.Background(), rm_well_link_sql, file.ID)
			if err != nil {
				return err
			}
			if fct.RowsAffected() < 1 {
				return fmt.Errorf("Prospect link not found for file id %d", file.ID)
			}
		}
		if len(file.WellIDs) > 0 {
			rm_well_link_sql, err := assets.ReadString("pg/file/delete_well_link_by_file_id.sql")
			if err != nil {
				return err
			}
			fct, err := tx.Exec(context.Background(), rm_well_link_sql, file.ID)
			if err != nil {
				return err
			}
			if fct.RowsAffected() < 1 {
				return fmt.Errorf("Well link not found for file id %d", file.ID)
			}
		}
	}
	rm_file_sql, err := assets.ReadString("pg/file/delete_by_file_id.sql")
	if err != nil {
		return err
	}
	fct, err := tx.Exec(context.Background(), rm_file_sql, file.ID)
	if err != nil {
		return err
	}
	if fct.RowsAffected() < 1 {
		return fmt.Errorf("File ID %d not found", file.ID)
	}
	// All files successfully added to the file table
	return tx.Commit(context.Background())
}
