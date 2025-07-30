package pg

import (
	"context"
	"fmt"
)

func (pg *Postgres) enumList(enum string) ([]string, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	clean_enum, err := conn.Conn().PgConn().EscapeString(enum)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("SELECT ARRAY_AGG(unnest ORDER BY unnest) "+
		"FROM UNNEST(ENUM_RANGE(null::\"%s\")::TEXT[])", clean_enum)
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	if rows.Next() {
		if err := rows.Scan(&values); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (pg *Postgres) enumAddValues(enum string, values ...string) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	e_enum, err := tx.Conn().PgConn().EscapeString(enum)
	if err != nil {
		return err
	}

	for _, value := range values {
		e_value, err := tx.Conn().PgConn().EscapeString(value)
		if err != nil {
			return err
		}
		q := fmt.Sprintf(
			"ALTER TYPE \"%s\" ADD VALUE IF NOT EXISTS '%s'", e_enum, e_value,
		)
		if _, err := tx.Exec(context.Background(), q); err != nil {
			return err
		}
	}
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) enumRename(enum, old_name, new_name string) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	clean_enum, err := tx.Conn().PgConn().EscapeString(enum)
	if err != nil {
		return err
	}
	clean_old_name, err := tx.Conn().PgConn().EscapeString(old_name)
	if err != nil {
		return err
	}
	clean_new_name, err := tx.Conn().PgConn().EscapeString(new_name)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(
		"ALTER TYPE \"%s\" RENAME VALUE '%s' TO '%s'", clean_enum, clean_old_name, clean_new_name,
	)
	if _, err := tx.Exec(context.Background(), q); err != nil {
		return err
	}
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
