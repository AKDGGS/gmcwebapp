package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) ListFilesById(tbl, id int) ([]*model.File, error) {
	sql := ""
	switch tbl {
	case model.BOREHOLE:
		sql = "pg/file/by_borehole_id.sql"
	case model.INVENTORY:
		sql = "pg/file/by_inventory_id.sql"
	case model.OUTCROP:
		sql = "pg/file/by_outcrop_id.sql"
	case model.PROSPECT:
		sql = "pg/file/by_prospect_id.sql"
	case model.WELL:
		sql = "pg/file/by_well_id.sql"
	}

	q, err := assets.ReadString(sql)
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	files, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.File])
	if err != nil {
		return nil, err
	}
	return files, nil
}
