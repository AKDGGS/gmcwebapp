package db

import (
	"context"
	"gmc/assets"
)

func (pg *Postgres) Drop() error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	dr, err := assets.ReadString("pg/misc/drop.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), dr); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
