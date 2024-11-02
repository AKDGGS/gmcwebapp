package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetInventoryByBarcode(barcode string, flags int) ([]*model.Inventory, error) {
	q, err := assets.ReadString("pg/inventory/by_barcode.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inventory, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.Inventory])
	if err != nil {
		return nil, err
	}
	return inventory, nil
}
