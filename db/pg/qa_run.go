package pg

import (
	"fmt"

	"gmc/db/model"
)

func (pg *Postgres) RunQAReport(id int) (*model.Table, error) {
	rpts, err := pg.ListQAReports()
	if err != nil {
		return nil, err
	}

	if id < 0 || id >= len(rpts) {
		return nil, fmt.Errorf("Invalid Report ID")
	}

	return pg.queryTable("pg/qa/" + rpts[id]["file"])
}
