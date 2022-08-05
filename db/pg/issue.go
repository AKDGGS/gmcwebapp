package pg

import (
	"context"
	"fmt"
	"strings"
)

func (pg *Postgres) ListIssues() ([]string, error) {
	return pg.enumList("issue")
}

func (pg *Postgres) AddIssues(issues ...string) error {
	return pg.enumAddValues("issue", issues...)
}

func (pg *Postgres) DeleteIssues(issues ...string) error {
	// Build new list of issues by removing the deleted issues
	// from the current list
	old_issues, err := pg.ListIssues()
	if err != nil {
		return err
	}
	new_issues := make(map[string]bool, len(old_issues))
	for _, ois := range old_issues {
		new_issues[ois] = true
	}
	for _, iss := range issues {
		delete(new_issues, iss)
	}

	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Step 1: Remove all references to the to-be-deleted issues
	del_sql := "UPDATE inventory_quality SET " +
		"issues = array_remove(issues, $1::issue) " +
		"WHERE issues @> ARRAY[$1::issue]"
	for _, iss := range issues {
		if _, err := tx.Exec(context.Background(), del_sql, iss); err != nil {
			return err
		}
	}

	// Step 2: Rename original issue type
	altertype_sql := "ALTER TYPE issue RENAME TO old_issue"
	if _, err := tx.Exec(context.Background(), altertype_sql); err != nil {
		return err
	}

	// Step 3a: Build statement to create new type
	var ntq strings.Builder
	for iss, _ := range new_issues {
		cis, err := tx.Conn().PgConn().EscapeString(iss)
		if err != nil {
			return err
		}

		if ntq.Len() > 0 {
			ntq.WriteString(", '")
		} else {
			ntq.WriteString("'")
		}
		ntq.WriteString(cis)
		ntq.WriteString("'")
	}

	// Step 3b: Create new type
	create_sql := fmt.Sprintf("CREATE TYPE issue AS ENUM(%s)", ntq.String())
	if _, err := tx.Exec(context.Background(), create_sql); err != nil {
		return err
	}

	// Step 4: Move inventory from old type to new type
	altercolumn_sql := "ALTER TABLE inventory_quality ALTER COLUMN issues " +
		"TYPE issue[] USING ((issues::text[])::issue[])"
	if _, err := tx.Exec(context.Background(), altercolumn_sql); err != nil {
		return err
	}

	// Step 5: Remove old type
	droptype_sql := "DROP TYPE old_issue"
	if _, err := tx.Exec(context.Background(), droptype_sql); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
