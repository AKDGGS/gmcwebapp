package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetFlatInventory(cb func(*model.FlatInventory) error) error {
	q, err := assets.ReadString("pg/search/inventory.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		f := model.FlatInventory{}
		err := rows.Scan(&f.ID, &f.Barcode, &f.Geometries)
		if err != nil {
			return err
		}
		if err := cb(&f); err != nil {
			return err
		}
	}
	return nil
}
