package pg

import (
	"context"
	"fmt"
)

func (pg *Postgres) AddKeywords(keywords ...string) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, kw := range keywords {
		k, err := tx.Conn().PgConn().EscapeString(kw)
		if err != nil {
			return err
		}
		q := fmt.Sprintf("ALTER TYPE keyword ADD VALUE IF NOT EXISTS '%s'", k)
		if _, err := tx.Exec(context.Background(), q); err != nil {
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
