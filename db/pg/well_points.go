package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetWellPoints() ([]model.WellPoint, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	q, err := assets.ReadString("pg/well/points.sql")
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), q)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var well_points []model.WellPoint
	for rows.Next() {
		var wp model.WellPoint
		err := rows.Scan(&wp.Name, &wp.WellID, &wp.Geog)
		if err != nil {
			return nil, err
		}
		well_points = append(well_points, wp)
	}
	return well_points, nil
}
