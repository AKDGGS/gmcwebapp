package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetFile(id int) (*model.File, error) {
	q, err := assets.ReadString("pg/file/by_file_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	file := model.File{}
	c, err := rowsToStruct(rows, &file)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, nil
	}
	return &file, nil
}
