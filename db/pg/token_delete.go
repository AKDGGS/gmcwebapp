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

	ct, err := pg.pool.Exec(context.Background(), q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() < 1 {
		return fmt.Errorf("token ID %d not found", id)
	}

	return nil
}
