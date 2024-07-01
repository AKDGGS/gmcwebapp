package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) CheckToken(tok string) (*model.Token, error) {
	q, err := assets.ReadString("pg/token/get_by_token.sql")
	if err != nil {
		return nil, err
	}

	r, err := pg.pool.Query(context.Background(), q, tok)
	if err != nil {
		return nil, err
	}

	tk, err := pgx.CollectOneRow(
		r, pgx.RowToAddrOfStructByNameLax[model.Token],
	)
	if err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return tk, nil
}
