package pg

import (
	"context"
	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) ListFiles(cb func(*model.File) error) error {
	q, err := assets.ReadString("pg/file/list.sql")
	if err != nil {
		return err
	}
	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		f, err := pgx.RowToAddrOfStructByNameLax[model.File](rows)
		if err != nil {
			return err
		}
		if err := cb(f); err != nil {
			return err
		}
	}
	return nil
}
