package pg

import (
	"fmt"
	"strings"

	dbf "gmc/db/flag"

	"github.com/jackc/pgtype"
)

func (pg *Postgres) GetInventory(id int, flags int) (map[string]interface{}, error) {
	inventory, err := pg.queryRow("pg/inventory/by_inventory_id.sql", id)
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

	t, ok := inventory["tray"].(int16)
	if !ok {
		delete(inventory, "tray")
	} else {
		inventory["tray"] = &t
	}

	if inventory["keywords"] != nil {
		kw, _ := inventory["keywords"]
		s := strings.Replace(fmt.Sprintf("%v", kw), " ", ", ", -1)
		if len(s) > 0 && s[len(s)-1] == ']' && s[0] == '[' {
			s = s[1 : len(s)-1]
		}
		inventory["keywords"] = &s
	}

	if (flags & dbf.FILES) != 0 {
		files, err := pg.queryRows("pg/file/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if files != nil {
			inventory["files"] = files
		}
	}

	if (flags & dbf.URLS) != 0 {
		urls, err := pg.queryRows("pg/url/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if urls != nil {
			inventory["urls"] = urls
		}
	}

	if (flags & dbf.NOTE) != 0 {
		notes, err := pg.queryRows("pg/note/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if notes != nil {
			inventory["notes"] = notes
		}
	}

	if (flags & dbf.PUBLICATION) != 0 {
		publications, err := pg.queryRows("pg/publication/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if publications != nil {
			inventory["publications"] = publications
		}
	}

	if (flags & dbf.BOREHOLE) != 0 {
		boreholes, err := pg.queryRows("pg/borehole/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if boreholes != nil {
			inventory["boreholes"] = boreholes
		}
	}

	if (flags & dbf.WELL) != 0 {
		wells, err := pg.queryRows("pg/well/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if wells != nil {
			inventory["wells"] = wells
		}
	}

	if (flags & dbf.SHOTLINE) != 0 {
		shotlines, err := pg.queryRows("pg/shotline/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if shotlines != nil {
			for _, m := range shotlines {
				for k, v := range m {
					if k == "shotpoint_number" {
						sp, ok := v.(pgtype.Numeric)
						if !ok {
							delete(inventory, "weight")
						} else {
							var ift float64
							sp.AssignTo(&ift)
							m["shotpoint_number"] = &ift
						}
					}
				}
			}
			inventory["shotlines"] = shotlines
		}
	}

	if (flags & dbf.OUTCROP) != 0 {
		outcrops, err := pg.queryRows("pg/outcrop/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if outcrops != nil {
			inventory["outcrops"] = outcrops
		}
	}

	if (flags & dbf.QUALITY) != 0 {
		qualities, err := pg.queryRows("pg/quality/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if qualities != nil {
			var issuesStr *string
			for _, m := range qualities {
				for k, v := range m {
					if k == "issues" {
						var s string
						if v != nil {
							s = strings.Replace(fmt.Sprintf("%v", v), " ", ", ", -1)
							s = strings.Replace(s, "_", " ", -1)
							if len(s) > 0 && s[len(s)-1] == ']' && s[0] == '[' {
								s = s[1 : len(s)-1]
							}
						} else {
							s = "GOOD"
						}
						issuesStr = &s
						if issuesStr != nil {
							m["issues"] = *issuesStr
						}
					}
				}
			}
			inventory["qualities"] = qualities
		}
	}

	if (flags & dbf.TRACKING) != 0 {
		containerlog, err := pg.queryRows("pg/container_log/by_inventory_id.sql", id)
		if err != nil {
			return nil, err
		}
		if containerlog != nil {
			inventory["containerlog"] = containerlog
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/inventory/geojson.sql", id)
		if err != nil {
			return nil, err
		}

		if geojson != nil {
			inventory["geojson"] = geojson["geojson"]
		}
	}

	return inventory, nil
}
