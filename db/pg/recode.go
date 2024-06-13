package pg

import (
	"context"
	"fmt"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) RecodeInventoryAndContainer(old_barcode string, new_barcode string) error {
	if old_barcode == "" || len(strings.TrimSpace(old_barcode)) < 1 {
		return fmt.Errorf("old barcode cannot be empty")
	}
	if new_barcode == "" || len(strings.TrimSpace(new_barcode)) < 1 {
		return fmt.Errorf("new barcode cannot be empty")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	// First, try updating the barcode
	q, err := assets.ReadString("pg/inventory/update_barcode.sql")
	if err != nil {
		return err
	}
	ir, err := tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/container/update_barcode.sql")
	if err != nil {
		return err
	}
	cr, err := tx.Exec(context.Background(), q, old_barcode, new_barcode)
	if err != nil {
		return err
	}
	rowsAffected := ir.RowsAffected() + cr.RowsAffected()
	// If those updates fail, try updating the alt_barcode
	if rowsAffected == 0 {
		q, err = assets.ReadString("pg/inventory/update_alt_barcode.sql")
		if err != nil {
			return err
		}
		ir, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
		if err != nil {
			return err
		}
		q, err = assets.ReadString("pg/container/update_alt_barcode.sql")
		if err != nil {
			return err
		}
		cr, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
		if err != nil {
			return err
		}
	}
	rowsAffected += ir.RowsAffected() + cr.RowsAffected()
	// If those updates fail, try appending "GMC-" to the barcode
	if (rowsAffected) == 0 {
		old_barcode = "GMC-" + old_barcode
		q, err = assets.ReadString("pg/inventory/update_barcode.sql")
		if err != nil {
			return err
		}
		ir, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
		if err != nil {
			return err
		}
		q, err = assets.ReadString("pg/container/update_barcode.sql")
		if err != nil {
			return err
		}
		cr, err = tx.Exec(context.Background(), q, old_barcode, new_barcode)
		if err != nil {
			return err
		}
	}
	rowsAffected += ir.RowsAffected() + cr.RowsAffected()
	if (rowsAffected) == 0 {
		return fmt.Errorf("recode failed")
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
