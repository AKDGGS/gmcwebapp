package pg

import (
	"context"
	"fmt"
	"strings"
)

func (pg *Postgres) ListContainerTypes() ([]string, error) {
	return pg.enumList("container_type")
}

func (pg *Postgres) AddContainerTypes(containertypes ...string) error {
	return pg.enumAddValues("container_type", containertypes...)
}

func (pg *Postgres) RenameContainerType(old_name, new_name string) error {
	return pg.enumRename("container_type", old_name, new_name)
}

func (pg *Postgres) DeleteContainerTypes(containertypes ...string) error {
	// Build new list of container types by removing the deleted container
	// types from the current list
	old, err := pg.ListContainerTypes()
	if err != nil {
		return err
	}
	new := make(map[string]bool, len(old))
	for _, oft := range old {
		new[oft] = true
	}
	for _, ft := range containertypes {
		delete(new, ft)
	}

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Step 1: Remove all references to the to-be-deleted types
	delquery := "UPDATE container SET " +
		"container_type = NULL " +
		"WHERE container_type = $1"
	for _, kw := range containertypes {
		if _, err := tx.Exec(context.Background(), delquery, kw); err != nil {
			return err
		}
	}

	// Step 2: Rename original type
	alterquery := "ALTER TYPE container_type RENAME TO old_container_type"
	if _, err := tx.Exec(context.Background(), alterquery); err != nil {
		return err
	}

	// Step 3a: Build statement to create new type
	var ntq strings.Builder
	for ct, _ := range new {
		ict, err := tx.Conn().PgConn().EscapeString(ct)
		if err != nil {
			return err
		}
		if ntq.Len() > 0 {
			ntq.WriteString(", ")
		}
		ntq.WriteString("'")
		ntq.WriteString(ict)
		ntq.WriteString("'")
	}

	// Step 3b: Create new type
	createquery := fmt.Sprintf("CREATE TYPE container_type AS ENUM(%s)", ntq.String())
	if _, err := tx.Exec(context.Background(), createquery); err != nil {
		return err
	}

	// Step 4: Move file from old type to new type
	altercolumnquery := "ALTER TABLE container ALTER COLUMN container_type " +
		"TYPE container_type USING ((container_type::text)::container_type)"
	if _, err := tx.Exec(context.Background(), altercolumnquery); err != nil {
		return err
	}

	// Step 5: Remove old type
	dropquery := "DROP TYPE old_container_type"
	if _, err := tx.Exec(context.Background(), dropquery); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
