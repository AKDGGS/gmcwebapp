package pg

import (
	"context"
	"strings"

	"gmc/assets"
	dbe "gmc/db/errors"
)

func (pg *Postgres) MoveInventoryAndContainersContents(src string, dest string) error {
	if src == "" || len(strings.TrimSpace(src)) < 1 {
		return dbe.ErrSourceBarcodeEmpty
	}
	if dest == "" || len(strings.TrimSpace(dest)) < 1 {
		return dbe.ErrDestinationBarcodeEmpty
	}
	q, err := assets.ReadString("pg/container/get_dest_container_id_and_validate_src_and_dest.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q, src, dest)
	if err != nil {
		return err
	}
	var dest_cid *int32
	var cid_count int32
	var src_valid bool
	for rows.Next() {
		if err := rows.Scan(&dest_cid, &cid_count, &src_valid); err != nil {
			return err
		}
	}
	if dest_cid == nil {
		return dbe.ErrDestinationNotFound
	}
	if !src_valid {
		return dbe.ErrSourceNotValid
	}
	if cid_count > 1 {
		return dbe.ErrDestinationMultipleContainers
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err = assets.ReadString("pg/container/move_by_parent_barcode.sql")
	if err != nil {
		return err
	}
	rc, err := tx.Exec(context.Background(), q, src, *dest_cid)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/inventory/move_by_container_barcode.sql")
	if err != nil {
		return err
	}
	ic, err := tx.Exec(context.Background(), q, src, dest_cid)
	if err != nil {
		return err
	}
	if ic.RowsAffected() == 0 {
		return dbe.ErrSrcNoInv
	}
	if (ic.RowsAffected() + rc.RowsAffected()) == 0 {
		return dbe.ErrNothingMoved
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
