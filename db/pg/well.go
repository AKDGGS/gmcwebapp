package pg

import (
	dbf "gmc/db/flag"

	"github.com/jackc/pgtype"
)

func (pg *Postgres) GetWell(id int, flags int) (map[string]interface{}, error) {
	well, err := pg.queryRow("pg/well_byid.sql", id)
	if err != nil {
		return nil, err
	}
	if well == nil {
		return nil, nil
	}

	md, ok := well["measured_depth"].(pgtype.Numeric)
	if !ok {
		delete(well, "measured_depth")
	} else {
		var ift float64
		md.AssignTo(&ift)
		well["measured_depth"] = &ift
	}

	vd, ok := well["vertical_depth"].(pgtype.Numeric)
	if !ok {
		delete(well, "vertical_depth")
	} else {
		var ift float64
		vd.AssignTo(&ift)
		well["vertical_depth"] = &ift
	}

	elv, ok := well["elevation"].(pgtype.Numeric)
	if !ok {
		delete(well, "elvation")
	} else {
		var ift float64
		elv.AssignTo(&ift)
		well["elevation"] = &ift
	}

	kb, ok := well["elevation_kb"].(pgtype.Numeric)
	if !ok {
		delete(well, "elevation_kb")
	} else {
		var ift float64
		kb.AssignTo(&ift)
		well["elevation_kb"] = &ift
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			well["files"] = files
		}
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		inventory, err := pg.queryRows(
			"pg/keyword_group_bywellid.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if inventory != nil {
			well["inventory"] = inventory
		}
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		operators, err := pg.queryRows("pg/organization_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if operators != nil {
			well["operators"] = operators
		}
	}

	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			well["urls"] = urls
		}
	}

	if (flags & dbf.NOTE) != 0 {
		notes, err := pg.queryRows("pg/note_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if notes != nil {
			well["notes"] = notes
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/well_geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson != nil {
			well["geojson"] = geojson["geojson"]
		}
	}

	if (flags & dbf.QUADRANGLES) != 0 {
		qds, err := pg.queryRows("pg/quadrangle250k_bywellid.sql", id)
		if err != nil {
			return nil, err
		}
		if qds != nil {
			well["quadrangles"] = qds
		}
	}

	return well, nil
}
