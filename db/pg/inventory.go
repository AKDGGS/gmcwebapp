package pg

import (
	"fmt"
	dbf "gmc/db/flag"
	"strings"

	"github.com/jackc/pgtype"
)

func (pg *Postgres) GetInventory(id int, flags int) (map[string]interface{}, error) {
	inventory, err := pg.queryRow("pg/inventory_byid.sql", id)
	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, nil
	}

	itop, ok := inventory["interval_top"].(pgtype.Numeric)
	if !ok {
		delete(inventory, "interval_top")
	} else {
		var ift float64
		itop.AssignTo(&ift)
		inventory["interval_top"] = &ift
	}

	ibot, ok := inventory["interval_bottom"].(pgtype.Numeric)
	if !ok {
		delete(inventory, "interval_bottom")
	} else {
		var ift float64
		ibot.AssignTo(&ift)
		inventory["interval_bottom"] = &ift
	}

	cd, ok := inventory["core_diameter"].(pgtype.Numeric)
	if !ok {
		delete(inventory, "core_diameter")
	} else {
		var ift float64
		cd.AssignTo(&ift)
		inventory["core_diameter"] = &ift
	}

	w, ok := inventory["weight"].(pgtype.Numeric)
	if !ok {
		delete(inventory, "weight")
	} else {
		var ift float64
		w.AssignTo(&ift)
		inventory["weight"] = &ift
	}

	kw, ok := inventory["keywords"].(pgtype.TextArray)
	if !ok {
		delete(inventory, "keywords")
	} else {
		var kws []string
		kw.AssignTo(&kws)
		s := strings.Join(kws[:], ", ")
		inventory["keywords"] = &s
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			inventory["files"] = files
		}
	}

	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			inventory["urls"] = urls
		}
	}

	if (flags & dbf.NOTE) != 0 {
		notes, err := pg.queryRows("pg/note_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		if notes != nil {
			inventory["notes"] = notes
		}
	}

	if (flags & dbf.RELATED) != 0 {
		type Any interface{}
		related := map[string]Any{}

		//Publications
		pubs, err := pg.queryRows("pg/publication_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		if pubs != nil {
			related["pubs"] = pubs

		}

		//Boreholes
		bores, err := pg.queryRows("pg/borehole_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		if bores != nil {
			related["bores"] = bores
		}

		//Wells
		wells, err := pg.queryRows("pg/well_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}

		if wells != nil {
			for k := range wells {
				well_id := wells[k]["well_id"]
				operators, errOrg := pg.queryRows("pg/organization_bywellid.sql", well_id)
				wells[k]["operators"] = operators
				if errOrg != nil {
					return nil, errOrg
				}
			}
			related["wells"] = wells
		}

		//Shotline
		shotlines, err := pg.queryRows("pg/shotline_byinventoryid.sql", id)
		if err != nil {
			return nil, err
		}
		for k := range shotlines {
			sp, ok := shotlines[k]["shotpoint_number"].(pgtype.Numeric)
			if !ok {
				delete(shotlines[k], "shotpoint_number")
			} else {
				var isp float64
				sp.AssignTo(&isp)
				shotlines[k]["shotpoint_number"] = &isp
			}
		}
		if shotlines != nil {
			related["shotlines"] = shotlines
		}
		inventory["related"] = related
	}
	fmt.Println(inventory)
	return inventory, nil
}
