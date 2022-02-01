package db

import (
	"context"
	"gmc/assets"
)

func (pg *Postgres) Init() error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Initialize types
	typ, err := assets.ReadString("pg/init/001-types.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), typ); err != nil {
		return err
	}

	// Initialize schema
	sch, err := assets.ReadString("pg/init/002-schema.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), sch); err != nil {
		return err
	}

	// Initialize views
	vw, err := assets.ReadString("pg/init/003-views.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), vw); err != nil {
		return err
	}

	// Initialize indexes
	idx, err := assets.ReadString("pg/init/004-indexes.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), idx); err != nil {
		return err
	}

	// Initialize triggers
	trg, err := assets.ReadString("pg/init/005-triggers.sql")
	if err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), trg); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
