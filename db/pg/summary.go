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

	count, err := cQryValue[int64](
		conn, "pg/container/get_count_by_barcode.sql", barcode,
	)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	q, err := assets.ReadString("pg/summary/by_barcode.sql")
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
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
