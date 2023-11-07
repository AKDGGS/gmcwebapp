package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) CheckToken(tok string) (*model.Token, error) {
	q, err := assets.ReadString("pg/token/get_by_token.sql")
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q, tok)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tk := model.Token{}

	c, err := rowsToStruct(rows, &tk)
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return nil, nil
	}
	return &tk, nil
}
