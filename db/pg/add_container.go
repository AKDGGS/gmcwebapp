package pg

import (
	"context"
	"fmt"

	"gmc/assets"
)

func (pg *Postgres) AddContainer(barcode string, name string, remark string) error {
	q, err := assets.ReadString("pg/container/get_count_by_barcode_inc_inventory.sql")
	if err != nil {
		return err
	}
	var count int
	err = pg.pool.QueryRow(context.Background(), q, barcode).Scan(&count)
	if err != nil {
		if err == ErrNoRows {
			return nil
		}
		return err
	}
	if count > 0 {
		return fmt.Errorf("barcode already exists")
	}
	q, err = assets.ReadString("pg/container/insert.sql")
	if err != nil {
		return err
	}
	_, err = pg.pool.Exec(context.Background(), q, barcode, name, remark)
	if err != nil {
		return err
	}
	return nil
}
