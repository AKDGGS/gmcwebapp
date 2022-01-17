package db

import (
	"time"
)

func (pg *Postgres) GetFile(id int, flags int) (int, string, time.Time, error) {
	/*
		var q string
		if include_private {
			q = assets.ReadString("pg/file_all.sql")
		} else {
			q = assets.ReadString("pg/file_public.sql")
		}

		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return 0, "", time.Time{}, err
		}
		defer rows.Close()

		if rows.Next() {
			var aid int
			var fname string
			var ftime time.Time
			err = rows.Scan(&aid, &fname, &ftime)
			if err != nil {
				return 0, "", time.Time{}, err
			}
			return aid, fname, ftime, nil
		} else {
			return 0, "", time.Time{}, nil
		}
	*/
	return 0, "", time.Time{}, nil
}
