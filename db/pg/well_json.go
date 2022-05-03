package pg

import (
	dbf "gmc/db/flag"
	"time"

	"github.com/jackc/pgtype"
)

func (pg *Postgres) GetWellJSON(id int, flags int) (map[string]interface{}, error) {
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

	cd, ok := well["completion_date"].(time.Time)
	if !ok {
		delete(well, "completion_date")
	} else {
		fcd := cd.Format("01-02-2006")
		well["completion_date"] = &fcd
	}

	sd, ok := well["spud_date"].(time.Time)
	if !ok {
		delete(well, "spud_date")
	} else {
		fsd := sd.Format("01-02-2006")
		well["spud_date"] = &fsd
	}

	if (flags & dbf.INVENTORY_SUMMARY) != 0 {
		keywords, err := pg.queryRows(
			"pg/keyword/group_by_well_id.sql", id,
			((flags & dbf.PRIVATE) == 0),
		)
		if err != nil {
			return nil, err
		}
		if keywords != nil {
			well["keywords"] = keywords
		}
	}

	return well, nil
}
