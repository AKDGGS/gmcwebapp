package pg

import (
	"context"

	"gmc/db/model"
)

func (pg *Postgres) ListFilesById(tbl, id int) ([]model.File, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	var files []model.File
	var sql = ""
	switch tbl {
	case model.BOREHOLE:
		sql = "pg/file/by_inventory_id.sql"
	case model.INVENTORY:
		sql = "pg/file/by_inventory_id.sql"
	case model.OUTCROP:
		sql = "pg/file/by_outcrop_id.sql"
	case model.PROSPECT:
		sql = "pg/file/by_prospect_id.sql"
	case model.WELL:
		sql = "pg/file/by_well_id.sql"
	}
	files, err = cQryStructs[model.File](conn, sql, id)
	if err != nil {
		return nil, err
	}
	return files, nil
}
