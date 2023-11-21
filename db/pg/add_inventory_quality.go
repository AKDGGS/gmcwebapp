package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) AddInventoryQuality(barcode string, remark string, issues []string, username string) error {
	if barcode == "" || len(strings.TrimSpace(barcode)) < 1 {
		return errors.New("Barcode cannot be empty")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err := assets.ReadString("pg/inventory/get_inventory_ids_by_barcode.sql")
	if err != nil {
		return err
	}
	var inventory_id int32
	err = tx.QueryRow(context.Background(), q, barcode).Scan(&inventory_id)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/quality/insert.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, inventory_id, remark, username, issues)
	if err != nil {
		return err
	}
	// If the insert is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
