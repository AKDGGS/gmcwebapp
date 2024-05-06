package pg

import (
	"context"
	"strings"

	"gmc/assets"
	dbe "gmc/db/errors"
)

func (pg *Postgres) AddInventory(barcode string, remark string, container_id *int32, issues []string, username string) error {
	if barcode == "" || len(strings.TrimSpace(barcode)) < 1 {
		return dbe.ErrBarcodeCannotBeEmpty
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err := assets.ReadString("pg/inventory/insert.sql")
	if err != nil {
		return err
	}
	var inventory_id int32
	err = tx.QueryRow(context.Background(), q, barcode, remark, nil).Scan(&inventory_id)
	if err != nil {
		return err
	}
	if inventory_id == 0 {
		return dbe.ErrInventoryInsertFailed
	}
	q, err = assets.ReadString("pg/quality/insert.sql")
	if err != nil {
		return err
	}
	var iq_id int32
	err = tx.QueryRow(context.Background(), q, inventory_id, nil, username, issues).Scan(&iq_id)
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
