package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetFile(id int) (*model.File, error) {
	q, err := assets.ReadString("pg/file/by_file_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	file := model.File{}
	c, err := rowsToStruct(rows, &file)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, nil
	}
	return &file, nil
}

func (pg *Postgres) PutFile(file *model.File, precommitFunc func() error) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	insert_sql, err := assets.ReadString("pg/file/insert.sql")
	if err != nil {
		return err
	}

	var file_id int32

	err = tx.QueryRow(context.Background(), insert_sql, file.Name, file.Description,
		file.Size, file.Type, file.MD5).Scan(&file_id)
	if err != nil {
		return err
	}

	file.ID = file_id

	if len(file.BoreholeIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_borehole.sql")
		if err != nil {
			return err
		}

		for _, b := range file.BoreholeIDs {
			if b != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file_id, b)
				if err != nil {
					return err
				}
			}
		}
	}
	if len(file.InventoryIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_inventory.sql")
		if err != nil {
			return err
		}

		for _, b := range file.InventoryIDs {
			if b != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file_id, b)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(file.OutcropIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_outcrop.sql")
		if err != nil {
			return err
		}

		for _, b := range file.OutcropIDs {
			if b != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file_id, b)
				if err != nil {
					return err
				}
			}
		}
	}
	if len(file.ProspectIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_prospect.sql")
		if err != nil {
			return err
		}

		for _, b := range file.ProspectIDs {
			if b != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file_id, b)
				if err != nil {
					return err
				}
			}
		}
	}
	if len(file.WellIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_well.sql")
		if err != nil {
			return err
		}

		for _, well := range file.WellIDs {
			if well != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file_id, well)
				if err != nil {
					return err
				}
			}
		}
	}

	if err = precommitFunc(); err != nil {
		return err
	}

	// All files successfully added to the file table
	tx.Commit(context.Background())
	return nil
}

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
	tx.Commit(context.Background())
	return nil
}
