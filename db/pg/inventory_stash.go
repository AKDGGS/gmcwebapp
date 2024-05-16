package pg

import (
	"context"

	"gmc/assets"
)

func (pg *Postgres) GetInventoryStash(id int) (interface{}, error) {
	sql, err := assets.ReadString("pg/inventory/stash.sql")
	if err != nil {
		return nil, err
	}

	var x interface{}
	row := pg.pool.QueryRow(context.Background(), sql, id)
	if err := row.Scan(&x); err != nil && err != ErrNoRows {
		return nil, err
	}
	return x, nil
}
