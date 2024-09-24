package pg

import (
	"context"
	"fmt"

	"gmc/assets"
)

func (pg *Postgres) MoveInventoryAndContainersContents(src string, dest string) error {
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
		return fmt.Errorf("destination barcode not found")
	}
	if !src_valid {
		return fmt.Errorf("source barcode not valid")
	}
	if cid_count > 1 {
		return fmt.Errorf("destination barcode refers to multiple containers")
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
	if _, err = tx.Exec(context.Background(), q, src, *dest_cid); err != nil {
		return err
	}
	q, err = assets.ReadString("pg/inventory/move_by_container_barcode.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), q, src, dest_cid); err != nil {
		return err
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
