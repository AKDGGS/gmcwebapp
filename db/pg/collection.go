package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) ListCollections() ([]*model.Collection, error) {
	q, err := assets.ReadString("pg/collection/list.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collections, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[model.Collection])
	if err != nil {
		return nil, err
	}
	return collections, nil
}
