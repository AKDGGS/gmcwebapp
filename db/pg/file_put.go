package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

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
	if len(file.Barcodes) > 0 {
		q, err := assets.ReadString("pg/inventory/get_ids_by_barcode.sql")
		if err != nil {
			return err
		}
		var inventory_ids []int32
		for _, barcode := range file.Barcodes {
			rows, err := tx.Query(context.Background(), q, barcode)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var inventory_id int32
				err := rows.Scan(&inventory_id)
				if err != nil {
					return err
				}
				inventory_ids = append(inventory_ids, inventory_id)
			}
		}
		if len(inventory_ids) == 0 {
			return fmt.Errorf("invalid barcode")
		}
		if len(inventory_ids) > 0 {
			insert_sql, err := assets.ReadString("pg/file/insert_inventory.sql")
			if err != nil {
				return err
			}
			for _, inv := range inventory_ids {
				if inv == 0 {
					continue
				}
				_, err = tx.Exec(context.Background(), insert_sql, file_id, inv)
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
