package pg

import (
	"context"

	"gmc/assets"
)

func (pg *Postgres) ListKeywords() ([]string, error) {
	q, err := assets.ReadString("pg/keyword/list.sql")
	if err != nil {
		return nil, err
	}

	rows, err := pg.pool.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []string
	if rows.Next() {
		if err := rows.Scan(&keywords); err != nil {
			return nil, err
		}
	}
	return keywords, nil
}
