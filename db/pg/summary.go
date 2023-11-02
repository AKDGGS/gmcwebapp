package pg

import (
	"context"
	"errors"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetSummaryByBarcode(barcode string, flags int) (*model.Summary, error) {
	summary := model.Summary{}
	q, err := assets.ReadString("pg/container/get_child_barcodes.sql")
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

	// return nil if the barcode is not a container
	if summary.Barcodes.BarcodeCount == 0 {
		return nil, errors.New("Barcode not found")
	}

	if (flags & dbf.CONTAINER_TOTAL) != 0 {
		q, err := assets.ReadString("pg/container/get_container_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Containers)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.COLLECTION_TOTAL) != 0 {
		q, err := assets.ReadString("pg/container/get_collection_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Collections)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.KEYWORD_SUMMARY) != 0 {
		q, err := assets.ReadString("pg/container/get_keyword_summary.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Keywords)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.BOREHOLE) != 0 {
		q, err := assets.ReadString("pg/container/get_borehole_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Boreholes)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.OUTCROP) != 0 {
		q, err := assets.ReadString("pg/container/get_outcrop_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Outcrops)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.SHOTLINE) != 0 {
		q, err := assets.ReadString("pg/container/get_shotline_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Shotlines)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.WELL) != 0 {
		q, err := assets.ReadString("pg/container/get_well_totals.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, barcode)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &summary.Wells)
		if err != nil {
			return nil, err
		}
	}
	return &summary, nil
}
