package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) MoveInventoryAndContainersContents(src string, dest string) error {
	if src == "" || len(strings.TrimSpace(src)) < 1 {
		return errors.New("Source barcode cannot be empty")
	}
	if dest == "" || len(strings.TrimSpace(dest)) < 1 {
		return errors.New("Destination barcode cannot be empty")
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
		return errors.New("The destination barcode does not exist")
	}
	if !src_valid {
		return errors.New("The source barcode is not valid")
	}
	if cid_count > 1 {
		return errors.New("The destination barcode refers to more than one inventory id")
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
	if (ic.RowsAffected() + rc.RowsAffected()) == 0 {
		return errors.New("The source barcode is empty")
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
