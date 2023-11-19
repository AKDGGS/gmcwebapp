package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) MoveInventoryAndContainers(dest string, barcodes_to_move []string, username string) error {
	if dest == "" || len(strings.TrimSpace(dest)) < 1 {
		return errors.New("Destination barcode cannot be empty")
	}
	if barcodes_to_move == nil || len(barcodes_to_move) < 1 {
		return errors.New("List of barcodes to be moved cannot be empty")
	}
	q, err := assets.ReadString("pg/container/get_container_id_by_barcode.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q, dest)
	if err != nil {
		return err
	}
	defer rows.Close()
	var dest_cid *int32
	var cid_count int32
	for rows.Next() {
		if err := rows.Scan(&dest_cid, &cid_count); err != nil {
			return err
		}
	}
	if cid_count != 1 {
		return errors.New("There was an problem with the destination barcode")
	}
	q, err = assets.ReadString("pg/move/validate_barcodes_to_move.sql")
	if err != nil {
		return err
	}
	var barcodes_valid bool
	err = pg.pool.QueryRow(context.Background(), q, barcodes_to_move).Scan(&barcodes_valid)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !barcodes_valid {
		return errors.New("At least one of the barcodes you are moving doesn't exist")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, barcode := range barcodes_to_move {
		q, err := assets.ReadString("pg/container/move_by_barcode.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), q, dest_cid, barcode)
		if err != nil {
			return err
		}
		q, err = assets.ReadString("pg/inventory/move_by_barcode.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), q, dest_cid, barcode)
		if err != nil {
			return err
		}
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
