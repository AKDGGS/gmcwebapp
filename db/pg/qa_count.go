package pg

import (
	"context"
	"fmt"

	"gmc/assets"
)

func (pg *Postgres) CountQAReport(id int) (int, error) {
	rpts, err := pg.ListQAReports()
	if err != nil {
		return 0, err
	}

	if id < 0 || id >= len(rpts) {
		return 0, fmt.Errorf("invalid Report ID")
	}

	q, err := assets.ReadString("pg/qa/" + rpts[id]["file"])
	if err != nil {
		return 0, err
	}
	q = "SELECT COUNT(*) FROM (" + q + ") AS q"

	var count int
	err = pg.pool.QueryRow(context.Background(), q).Scan(&count)
	if err != nil {
		if err == ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
