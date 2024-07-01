package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetFile(id int) (*model.File, error) {
	q, err := assets.ReadString("pg/file/by_file_id.sql")
	if err != nil {
		return nil, err
	}
	r, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	f, err := pgx.CollectOneRow(
		r, pgx.RowToAddrOfStructByNameLax[model.File],
	)
	if err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}
