package pg

import (
	"context"

	"gmc/assets"
)

func (pg *Postgres) GetFileNextID() (int32, error) {
	q, err := assets.ReadString("pg/file/next_id.sql")
	if err != nil {
		return 0, err
	}

	var id int32
	err = pg.pool.QueryRow(context.Background(), q).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
