package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) Audit(barcode string) ([]*model.ContainerBarcodes, error) {
	get_cid_query, err := assets.ReadString("pg/audit/get_container_id_by_barcode.sql")
	if err != nil {
		return nil, err
	}
	var container_id int
	err = pg.pool.QueryRow(context.Background(), get_cid_query, barcode).Scan(&container_id)
	if err != nil {
		if err == ErrNoRows {
			return nil, fmt.Errorf("barcode not found")
		}
		return nil, err
	}
	get_barcodes_query, err := assets.ReadString("pg/audit/get_barcodes_by_container_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), get_barcodes_query, container_id)
	if err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.ContainerBarcodes])
	if err != nil {
		return nil, err
	}
	return results, nil
}
