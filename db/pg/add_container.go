package pg

import (
	"context"
	"errors"
	"strings"

	"gmc/assets"
	authu "gmc/auth/util"
)

func (pg *Postgres) AddContainer(barcode string, alt_barcode string, name string, remark string, user *authu.User) error {
	if nil == user {
		return errors.New("Access denied.")
	}
	if barcode == "" || len(strings.TrimSpace(barcode)) < 1 {
		return errors.New("Barcode cannot be empty")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	q, err := assets.ReadString("pg/container/container_type_by_name.sql")
	if err != nil {
		return err
	}
	var container_type_id int32
	err = tx.QueryRow(context.Background(), q, "unknown").Scan(&container_type_id)
	if err != nil {
		return err
	}
	q, err = assets.ReadString("pg/container/insert.sql")
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), q, container_type_id, barcode, nil, name, remark)
	if err != nil {
		return err
	}
	// If the insert is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
