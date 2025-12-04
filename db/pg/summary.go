package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetSummaryByBarcode(barcode string, flags int) (*model.Summary, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	q, err := assets.ReadString("pg/summary/by_barcode.sql")
	if err == ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary, err := pgx.CollectOneRow(
		rows, pgx.RowToStructByNameLax[model.Summary],
	)
	if err != nil && err != ErrNoRows {
		return nil, err
	}
	return &summary, nil
}
