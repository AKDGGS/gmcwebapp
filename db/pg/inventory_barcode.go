package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetInventoryByBarcode(barcode string, flags int) ([]*model.Inventory, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	q, err := assets.ReadString("pg/inventory/by_barcode.sql")
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inventory, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Inventory])
	if err != nil {
		return nil, err
	}
	result := make([]*model.Inventory, len(inventory))
	for i, inv := range inventory {
		result[i] = &inv
	}
	return result, nil
}
