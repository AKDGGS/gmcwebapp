package pg

import (
	"context"
	"strings"

	"gmc/assets"
	dbe "gmc/db/errors"
)

func (pg *Postgres) AddInventoryQuality(barcode string, remark string, issues []string, username string) error {
	if barcode == "" || len(strings.TrimSpace(barcode)) < 1 {
		return dbe.ErrSourceBarcodeEmpty
	}
	q, err := assets.ReadString("pg/inventory/get_ids_by_barcode.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return err
	}
	defer rows.Close()

	var inventory_ids []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return err
		}
		inventory_ids = append(inventory_ids, id)
	}
	if inventory_ids == nil {
		return dbe.ErrNotFoundInInventory
	}
	if len(inventory_ids) > 1 {
		return dbe.ErrMultipleIDs
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err = assets.ReadString("pg/quality/insert.sql")
	if err != nil {
		return err
	}
	var iq_id int32
	err = tx.QueryRow(context.Background(), q, inventory_ids[0], nil, username, issues).Scan(&iq_id)
	if err != nil {
		return err
	}
	if iq_id == 0 {
		return dbe.ErrInventoryQualityInsertFailed
	}
	// If the insert is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
