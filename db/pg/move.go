package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
	authu "gmc/auth/util"
)

func (pg *Postgres) MoveByBarcode(dest string, container_list []string, user *authu.User) error {
	if nil == user {
		return errors.New("Access denied.")
	}
	if dest == "" || len(strings.TrimSpace(dest)) < 1 {
		return errors.New("Destination barcode cannot be empty")
	}
	if container_list == nil || len(container_list) < 1 {
		return errors.New("List of barcodes to be moved cannot be empty.")
	}
	q, err := assets.ReadString("pg/container/get_container_ids_by_barcode.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q, dest)
	if err != nil {
		return err
	}
	defer rows.Close()
	var dest_container_id int32
	for rows.Next() {
		if err := rows.Scan(&dest_container_id); err != nil {
			return err
		}
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, barcode := range container_list {
		q, err := assets.ReadString("pg/move/container_by_barcode.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), q, dest_container_id, barcode)
		if err != nil {
			return err
		}
		q, err = assets.ReadString("pg/move/inventory_by_barcode.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), q, dest_container_id, barcode)
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
