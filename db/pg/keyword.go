package pg

import (
	"context"
	"fmt"
	"strings"
)

func (pg *Postgres) ListKeywords() ([]string, error) {
	return pg.enumList("keyword")
}

func (pg *Postgres) AddKeywords(keywords ...string) error {
	return pg.enumAddValues("keyword", keywords...)
}

func (pg *Postgres) RenameKeyword(old_name, new_name string) error {
	return pg.enumRename("keyword", old_name, new_name)
}

func (pg *Postgres) DeleteKeywords(keywords ...string) error {
	// Build new list of keywords by removing the deleted keywords
	// from the current list
	old_keywords, err := pg.ListKeywords()
	if err != nil {
		return err
	}
	new_keywords := make(map[string]bool, len(old_keywords))
	for _, okw := range old_keywords {
		new_keywords[okw] = true
	}
	for _, kw := range keywords {
		delete(new_keywords, kw)
	}

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Step 1: Remove all references to the to-be-deleted keys
	del_sql := "UPDATE inventory SET " +
		"keywords = array_remove(keywords, $1::keyword) " +
		"WHERE keywords @> ARRAY[$1::keyword]"
	for _, kw := range keywords {
		if _, err := tx.Exec(context.Background(), del_sql, kw); err != nil {
			return err
		}
	}

	// Step 2: Rename original keyword type
	altertype_sql := "ALTER TYPE keyword RENAME TO old_keyword"
	if _, err := tx.Exec(context.Background(), altertype_sql); err != nil {
		return err
	}

	// Step 3a: Build statement to create new type
	var ntq strings.Builder
	for k, _ := range new_keywords {
		ik, err := tx.Conn().PgConn().EscapeString(k)
		if err != nil {
			return err
		}
		if ntq.Len() > 0 {
			ntq.WriteString(", ")
		}
		ntq.WriteString("'")
		ntq.WriteString(ik)
		ntq.WriteString("'")
	}

	// Step 3b: Create new type
	create_sql := fmt.Sprintf("CREATE TYPE keyword AS ENUM(%s)", ntq.String())
	if _, err := tx.Exec(context.Background(), create_sql); err != nil {
		return err
	}

	// Step 4: Move inventory from old type to new type
	altercolumn_sql := "ALTER TABLE inventory ALTER COLUMN keywords " +
		"TYPE keyword[] USING ((keywords::text[])::keyword[])"
	if _, err := tx.Exec(context.Background(), altercolumn_sql); err != nil {
		return err
	}

	// Step 5: Remove old type
	droptype_sql := "DROP TYPE old_keyword"
	if _, err := tx.Exec(context.Background(), droptype_sql); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
