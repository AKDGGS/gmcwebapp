package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) CreateToken(tk *model.Token) error {
	q, err := assets.ReadString("pg/token/create.sql")
	if err != nil {
		return err
	}

	rows, err := pg.pool.Query(context.Background(), q, tk.Description, tk.Token)
	if err != nil {
		return err
	}
	if rows.Next() {
		err = rows.Scan(&tk.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
