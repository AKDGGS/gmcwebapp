package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) MoveContents(src string, dest string) error {
	if src == "" || len(strings.TrimSpace(src)) < 1 {
		return errors.New("Source barcode cannot be empty")
	}
	if dest == "" || len(strings.TrimSpace(dest)) < 1 {
		return errors.New("Destination barcode cannot be empty")
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
	var container_ids []int32
	for rows.Next() {
		var container_id int32
		if err := rows.Scan(&container_id); err != nil {
			return err
		}
		container_ids = append(container_ids, container_id)
	}
	if len(container_ids) > 1 {
		return errors.New("Destination barcode refers to multiple containers")
	}
	container_id := container_ids[0]
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err = assets.ReadString("pg/move/container_by_parent_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, src, container_id)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/move/inventory_by_container_barcode.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, src, container_id)
	if err != nil {
		return err
	}
	// If the move is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
