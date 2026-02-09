package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) PutFile(file *model.File, precommitFunc func(*model.File) error) error {
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

	insert_sql, err := assets.ReadString("pg/file/insert.sql")
	if err != nil {
		return err
	}
	var file_id int32
	err = tx.QueryRow(context.Background(), insert_sql, file.Name, file.Description,
		file.Size, file.Type).Scan(&file_id)
	if err != nil {
		if err == ErrNoRows {
			return nil
		}
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

	if err = precommitFunc(file); err != nil {
		return err
	}

	if file.Hash != nil {
		hashsql, err := assets.ReadString("pg/file/update_hash.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), hashsql, file.Hash, file.ID)
		if err != nil {
			return err
		}
		fmt.Println(file.Hash)
	}

	// All files successfully added to the file table
	return tx.Commit(context.Background())
}
