package pg

import (
	"context"
	"fmt"
	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetOutcrop(id int, flags int) (*model.Outcrop, error) {
	q, err := assets.ReadString("pg/outcrop/by_outcrop_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oc := model.Outcrop{}
	rowToStruct(rows, &oc)
	fmt.Println(oc)

	//
	// if (flags & dbf.FILES) != 0 {
	// 	files, err := pg.queryRows("pg/file/by_outcrop_id.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if files != nil {
	// 		outcrop["files"] = files
	// 	}
	// }
	//
	// if (flags & dbf.INVENTORY_SUMMARY) != 0 {
	// 	kw, err := pg.queryRows(
	// 		"pg/keyword/group_by_outcrop_id.sql", id,
	// 		((flags & dbf.PRIVATE) == 0),
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if kw != nil {
	// 		outcrop["keywords"] = kw
	// 	}
	// }
	//
	// if (flags & dbf.ORGANIZATION) != 0 {
	// 	organizations, err := pg.queryRows("pg/organization/by_outcrop_id.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if organizations != nil {
	// 		outcrop["organizations"] = organizations
	// 	}
	// }
	//
	// if (flags & dbf.URLS) != 0 {
	// 	urls, err := pg.queryRows("pg/url/by_outcrop_id.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if urls != nil {
	// 		outcrop["urls"] = urls
	// 	}
	// }
	//
	// if (flags & dbf.NOTE) != 0 {
	// 	notes, err := pg.queryRows("pg/note/by_outcrop_id.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if notes != nil {
	// 		outcrop["notes"] = notes
	// 	}
	// }
	//
	// if (flags & dbf.GEOJSON) != 0 {
	// 	geojson, err := pg.queryRow("pg/outcrop/geojson.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if geojson != nil {
	// 		outcrop["geojson"] = geojson["geojson"]
	// 	}
	// }
	//
	// if (flags & dbf.QUADRANGLES) != 0 {
	// 	qds, err := pg.queryRows("pg/quadrangle/250k_by_outcrop_id.sql", id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if qds != nil {
	// 		outcrop["quadrangles"] = qds
	// 	}
	// }

	return &oc, nil
}
