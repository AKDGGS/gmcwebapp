package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) ListTokens() ([]*model.Token, error) {
	q, err := assets.ReadString("pg/token/list.sql")
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]*model.Token, 0)
	for rows.Next() {
		var tk model.Token
		err = rows.Scan(&tk.ID, &tk.Description)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &tk)
	}
	return tokens, nil
}
