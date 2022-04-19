package pg

import (
	dbf "gmc/db/flag"

	"github.com/jackc/pgtype"
)

func (pg *Postgres) GetShotline(id int, flags int) (map[string]interface{}, error) {
	shotline, err := pg.queryRow("pg/shotline_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if shotline == nil {
		return nil, nil
	}

	ptmin, ok := shotline["shotpoint_min"].(pgtype.Numeric)
	if !ok {
		delete(shotline, "shotpoint_min")
	} else {
		var ift float64
		ptmin.AssignTo(&ift)
		shotline["shotpoint_min"] = &ift
	}

	ptmax, ok := shotline["shotpoint_max"].(pgtype.Numeric)
	if !ok {
		delete(shotline, "shotpoint_max")
	} else {
		var ift float64
		ptmax.AssignTo(&ift)
		shotline["shotpoint_max"] = &ift
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword/group_by_shotline_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			shotline["inventory"] = inventory
		}
	}
	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url/by_shotline_id.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			shotline["urls"] = urls
		}
	}

	if (flags & dbf.NOTE) != 0 {
		notes, err := pg.queryRows("pg/note/by_shotline_id.sql", id)
		if err != nil {
			return nil, err
		}
		if notes != nil {
			shotline["notes"] = notes
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/shotline_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			shotline["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle/250k_by_shotline_id.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			shotline["quadrangles"] = qds
		}
	}

	return shotline, nil
}
