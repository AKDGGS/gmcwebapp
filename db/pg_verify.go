package db

import (
	"context"
	"fmt"
)

func (pg *Postgres) Verify() error {
	row, err := pg.pool.Query(context.Background(), "SELECT 1")
	if err != nil {
		return err
	}
	defer row.Close()

	if !row.Next() {
		return fmt.Errorf("SELECT 1 returned no rows")
	}

	var one int
	if err := row.Scan(&one); err != nil {
		return err
	}

	if one != 1 {
		return fmt.Errorf("SELECT 1 returned value not equal to 1")
	}
	return nil
}
