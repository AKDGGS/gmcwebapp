package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetSummaryByBarcode(barcode string, flags int) (*model.Summary, error) {
	q, err := assets.ReadString("pg/container/get_count_by_barcode.sql")
	if err != nil {
		return nil, err
	}
	var barcode_count int32
	err = pg.pool.QueryRow(context.Background(), q, barcode).Scan(&barcode_count)
	if err != nil {
		return nil, err
	}
	// return nil if the barcode is not a container (barcode_count == 0)
	if barcode_count == 0 {
		return nil, nil
	}

	summary := model.Summary{}
	q, err = assets.ReadString("pg/container/get_child_barcodes.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Barcodes)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Containers)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_collection_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Collections)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_keyword_summary.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Keywords)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_borehole_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Boreholes)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_outcrop_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Outcrops)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_shotline_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Shotlines)
	if err != nil {
		return nil, err
	}

	q, err = assets.ReadString("pg/container/get_well_totals.sql")
	if err != nil {
		return nil, err
	}
	rows, err = pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	_, err = rowsToStruct(rows, &summary.Wells)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
