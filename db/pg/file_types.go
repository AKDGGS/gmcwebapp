package pg

import (
	"context"
	"fmt"
	"strings"
)

func (pg *Postgres) ListFileTypes() ([]string, error) {
	return pg.enumList("file_types")
}

func (pg *Postgres) AddFileTypes(filetypes ...string) error {
	return pg.enumAddValues("file_types", filetypes...)
}

func (pg *Postgres) RenameFileTypes(old, new string) error {
	return pg.enumRename("file_types", old, new)
}

func (pg *Postgres) DeleteFileTypes(filetypes ...string) error {
	// Build new list of file types by removing the deleted file types
	// from the current list
	old, err := pg.ListFileTypes()
	if err != nil {
		return err
	}
	new := make(map[string]bool, len(old))
	for _, oft := range old {
		new[oft] = true
	}
	for _, ft := range filetypes {
		delete(new, ft)
	}

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Step 1: Remove all references to the to-be-deleted types
	delquery := "UPDATE file SET " +
		"file_type = NULL " +
		"WHERE file_type = $1"
	for _, kw := range filetypes {
		if _, err := tx.Exec(context.Background(), delquery, kw); err != nil {
			return err
		}
	}

	// Step 2: Rename original type
	alterquery := "ALTER TYPE file_types RENAME TO old_filetype"
	if _, err := tx.Exec(context.Background(), alterquery); err != nil {
		return err
	}

	// Step 3a: Build statement to create new type
	var ntq strings.Builder
	for ft, _ := range new {
		ift, err := tx.Conn().PgConn().EscapeString(ft)
		if err != nil {
			return err
		}
		if ntq.Len() > 0 {
			ntq.WriteString(", ")
		}
		ntq.WriteString("'")
		ntq.WriteString(ift)
		ntq.WriteString("'")
	}

	// Step 3b: Create new type
	createquery := fmt.Sprintf("CREATE TYPE file_types AS ENUM(%s)", ntq.String())
	if _, err := tx.Exec(context.Background(), createquery); err != nil {
		return err
	}

	// Step 4: Move file from old type to new type
	altercolumnquery := "ALTER TABLE file ALTER COLUMN file_type " +
		"TYPE file_types USING ((file_type::text)::file_types)"
	if _, err := tx.Exec(context.Background(), altercolumnquery); err != nil {
		return err
	}

	// Step 5: Remove old type
	dropquery := "DROP TYPE old_filetype"
	if _, err := tx.Exec(context.Background(), dropquery); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
