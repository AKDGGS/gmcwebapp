package pg

import (
	"context"

	"gmc/assets"
)

func (pg *Postgres) GetWellPoints() (interface{}, error) {
	sql, err := assets.ReadString("pg/well/points.sql")
	if err != nil {
		return nil, err
	}

	var x interface{}
	row := pg.pool.QueryRow(context.Background(), sql)
	if err := row.Scan(&x); err != nil && err != ErrNoRows {
		return nil, err
	}
	return x, nil
}
