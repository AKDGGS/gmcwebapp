package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) PutFile(file *model.File, precommitFunc func(*model.File) (bool, error)) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if file.Barcode != "" {
		q, err := assets.ReadString("pg/inventory/get_ids_by_barcode.sql")
		if err != nil {
			return err
		}
		rows, err := tx.Query(context.Background(), q, file.Barcode)
		if err != nil {
			return err
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				return err
			}
			file.InventoryIDs = append(file.InventoryIDs, id)
			count++
		}

		if count < 1 {
			return fmt.Errorf(
				"no inventory found matching barcode \"%s\"", file.Barcode,
			)
		}
	}

	if file.ID == 0 {
		// Handle an insert
		insert_sql, err := assets.ReadString("pg/file/insert.sql")
		if err != nil {
			return err
		}
		err = tx.QueryRow(context.Background(), insert_sql,
			file.Name, file.Description, file.Size, file.Type,
		).Scan(&file.ID)
		if err != nil {
			if err == ErrNoRows {
				return nil
			}
			return err
		}
	} else {
		// Handle an update
		update_sql, err := assets.ReadString("pg/file/update.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			context.Background(), update_sql,
			file.Name, file.Description, file.Size,
			file.Type, file.Hash, file.ID,
		)
		if err != nil {
			return err
		}
	}

	if len(file.InventoryIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_inventory.sql")
		if err != nil {
			return err
		}
		for _, id := range file.InventoryIDs {
			if id != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file.ID, id)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(file.BoreholeIDs) > 0 {
		insert_sql, err := assets.ReadString("pg/file/insert_borehole.sql")
		if err != nil {
			return err
		}
		for _, id := range file.BoreholeIDs {
			if id != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file.ID, id)
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
		for _, id := range file.OutcropIDs {
			if id != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file.ID, id)
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
		for _, id := range file.ProspectIDs {
			if id != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file.ID, id)
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
		for _, id := range file.WellIDs {
			if id != 0 {
				_, err = tx.Exec(context.Background(), insert_sql, file.ID, id)
				if err != nil {
					return err
				}
			}
		}
	}

	if precommitFunc != nil {
		req_update, err := precommitFunc(file)
		if err != nil {
			return err
		}
		if req_update {
			update_sql, err := assets.ReadString("pg/file/update.sql")
			if err != nil {
				return err
			}
			_, err = tx.Exec(
				context.Background(), update_sql,
				file.Name, file.Description, file.Size,
				file.Type, file.Hash, file.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(context.Background())
}
