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

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(), q, tk.Description, tk.Token,
	).Scan(&tk.ID)
	if err != nil {
		if err == ErrNoRows {
			return nil
		}
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
