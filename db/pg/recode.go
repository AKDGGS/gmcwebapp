package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) RecodeInventoryAndContainer(old_barcode string, new_barcode string) error {
	if old_barcode == "" || len(strings.TrimSpace(old_barcode)) < 1 {
		return errors.New("Old barcode cannot be empty")
	}
	if new_barcode == "" || len(strings.TrimSpace(new_barcode)) < 1 {
		return errors.New("New barcode cannot be empty")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err := assets.ReadString("pg/container/update_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/container/update_alt_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/inventory/update_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/inventory/update_alt_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
