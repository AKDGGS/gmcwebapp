package pg

import (
	"context"
	"fmt"

	"gmc/assets"
)

func (pg *Postgres) DeleteToken(id int) error {
	q, err := assets.ReadString("pg/token/delete.sql")
	if err != nil {
		return err
	}

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	ct, err := tx.Exec(context.Background(), q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() < 1 {
		return fmt.Errorf("token ID %d not found", id)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
