package pg

import (
	"context"
	"errors"

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

	if rowsToStruct(rows, &tk) == 0 {
		return nil, errors.New("No matches found")
	}
	return &tk, nil
}
