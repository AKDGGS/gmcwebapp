package pg

import (
	"context"
	"strings"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetInventory(id int, flags int) (*model.Inventory, error) {
	q, err := assets.ReadString("pg/inventory/by_inventory_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inventory := model.Inventory{}

	if count, err := rowsToStruct(rows, &inventory); err != nil || count == 0 {
		return nil, err
	}

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Files)
	}

	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(r, &inventory.URLs)
	}

	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Notes)
	}

	if (flags & dbf.PUBLICATION) != 0 {
		q, err = assets.ReadString("pg/publication/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Publications)
	}

	if (flags & dbf.BOREHOLE) != 0 {
		q, err := assets.ReadString("pg/borehole/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Boreholes)
	}

	if (flags & dbf.OUTCROP) != 0 {
		q, err := assets.ReadString("pg/outcrop/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Outcrops)
	}

	if (flags & dbf.SHOTPOINT) != 0 {
		q, err := assets.ReadString("pg/shotpoint/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Shotpoints)
	}

	if (flags & dbf.WELL) != 0 {
		q, err := assets.ReadString("pg/well/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Wells)
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		for i := 0; i < len(inventory.Boreholes); i++ {
			q, err := assets.ReadString("pg/organization/by_borehole_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory.Boreholes[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowsToStruct(rows, &inventory.Boreholes[i].Organizations)
		}
		for i := 0; i < len(inventory.Outcrops); i++ {
			q, err := assets.ReadString("pg/organization/by_outcrop_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory.Outcrops[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowsToStruct(rows, &inventory.Outcrops[i].Organizations)
		}
		for i := 0; i < len(inventory.Wells); i++ {
			q, err := assets.ReadString("pg/organization/by_well_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory.Wells[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowsToStruct(rows, &inventory.Wells[i].Organizations)
		}
	}

	if (flags & dbf.QUALITY) != 0 {
		q, err := assets.ReadString("pg/quality/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.Qualities)
	}

	if (flags & dbf.TRACKING) != 0 {
		q, err := assets.ReadString("pg/container_log/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowsToStruct(rows, &inventory.ContainerLog)
	}

	if (flags & dbf.GEOJSON) != 0 {
		geojson, err := pg.queryRow("pg/inventory/geojson.sql", id)
		if err != nil {
			return nil, err
		}
		if geojson["geojson"] != nil {
			inventory.GeoJSON = geojson["geojson"].(map[string]interface{})
		}
	}
	return &inventory, nil
}

func (pg *Postgres) GetInventoryByBarcode(barcode string, flags int) ([]*model.Inventory, error) {
	q, err := assets.ReadString("pg/inventory/by_barcode.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := []*model.Inventory{}
	bh_inv := make(map[int32][]*model.Inventory)
	o_inv := make(map[int32][]*model.Inventory)
	sp_inv := make(map[int32][]*model.Inventory)
	w_inv := make(map[int32][]*model.Inventory)
	q_inv := make(map[int32][]*model.Inventory)
	var borehole_ids, outcrop_ids, shotpoint_ids, well_ids, quality_ids []int32
	var t_bhids, t_oids, t_spids, t_wids, t_qids []int32
	cols := rows.FieldDescriptions()
	ptrs_list := make([]interface{}, len(cols))

	for i := 0; i < len(cols); i++ {
		switch strings.ToLower(string(cols[i].Name)) {
		case "borehole_ids":
			ptrs_list[i] = &t_bhids
		case "outcrop_ids":
			ptrs_list[i] = &t_oids
		case "shotpoint_ids":
			ptrs_list[i] = &t_spids
		case "well_ids":
			ptrs_list[i] = &t_wids
		case "quality_ids":
			ptrs_list[i] = &t_qids
		}
	}

	for rows.Next() {
		if err := rows.Scan(ptrs_list...); err != nil {
			return nil, err
		}
		inventory = append(inventory, &model.Inventory{})
		rowToStruct(rows, inventory[len(inventory)-1])
		if len(t_bhids) > 0 {
			for _, borehole_id := range t_bhids {
				bh_inv[borehole_id] = append(bh_inv[borehole_id], inventory[len(inventory)-1])
			}
			borehole_ids = append(borehole_ids, t_bhids...)
		}
		if len(t_oids) > 0 {
			for _, outcrop_id := range t_oids {
				o_inv[outcrop_id] = append(o_inv[outcrop_id], inventory[len(inventory)-1])
			}
			outcrop_ids = append(outcrop_ids, t_oids...)
		}
		if len(t_spids) > 0 {
			for _, shotpoint_id := range t_spids {
				sp_inv[shotpoint_id] = append(sp_inv[shotpoint_id], inventory[len(inventory)-1])
			}
			shotpoint_ids = append(shotpoint_ids, t_spids...)
		}
		if len(t_wids) > 0 {
			for _, well_id := range t_wids {
				w_inv[well_id] = append(w_inv[well_id], inventory[len(inventory)-1])
			}
			well_ids = append(well_ids, t_wids...)
		}
		qids_exists := make(map[int32]bool)
		if len(t_qids) > 0 {
			for _, quality_id := range t_qids {
				if _, exists := qids_exists[quality_id]; !exists {
					qids_exists[quality_id] = true
					q_inv[quality_id] = append(q_inv[quality_id], inventory[len(inventory)-1])
					quality_ids = append(quality_ids, quality_id)
				}
			}
		}
	}

	if len(borehole_ids) > 0 {
		q, err := assets.ReadString("pg/borehole/by_borehole_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, borehole_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		boreholes := make([]model.Borehole, 0)
		rowsToStruct(rows, &boreholes)
		for _, borehole := range boreholes {
			for _, inv := range bh_inv[borehole.ID] {
				inv.Boreholes = append(inv.Boreholes, borehole)
			}
		}
	}

	if len(outcrop_ids) > 0 {
		q, err := assets.ReadString("pg/outcrop/by_outcrop_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, outcrop_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		outcrops := make([]model.Outcrop, 0)
		rowsToStruct(rows, &outcrops)
		for _, outcrop := range outcrops {
			for _, inv := range o_inv[outcrop.ID] {
				inv.Outcrops = append(inv.Outcrops, outcrop)
			}
		}
	}

	if len(shotpoint_ids) > 0 {
		q, err := assets.ReadString("pg/shotpoint/by_shotpoint_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, shotpoint_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		shotpoints := make([]model.Shotpoint, 0)
		rowsToStruct(rows, &shotpoints)
		for _, sp := range shotpoints {
			for _, inv := range sp_inv[sp.ID] {
				inv.Shotpoints = append(inv.Shotpoints, sp)
			}
		}
	}

	if len(well_ids) > 0 {
		q, err := assets.ReadString("pg/well/by_well_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, well_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		o_well := make(map[int32][]*model.Well)
		wells_cols := rows.FieldDescriptions()
		wptrs_list := make([]interface{}, len(wells_cols))
		var t_orgids []int32
		for k := 0; k < len(wells_cols); k++ {
			if strings.EqualFold(string(wells_cols[k].Name), "operator_ids") {
				wptrs_list[k] = &t_orgids
			}
		}

		wells := make([]*model.Well, 0)
		for rows.Next() {
			if err := rows.Scan(wptrs_list...); err != nil {
				return nil, err
			}
			wells = append(wells, &model.Well{})
			rowToStruct(rows, wells[len(wells)-1])
			if len(t_orgids) > 0 {
				for _, org_id := range t_orgids {
					o_well[org_id] = append(o_well[org_id], wells[len(wells)-1])
				}
			}
		}

		q, err = assets.ReadString("pg/organization/by_well_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err = pg.pool.Query(context.Background(), q, well_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		orgs := make([]model.Organization, 0)
		rowsToStruct(rows, &orgs)
		for _, org := range orgs {
			for _, w := range o_well[org.ID] {
				w.Organizations = append(w.Organizations, org)
			}
		}
		for _, w := range wells {
			for _, inv := range w_inv[w.ID] {
				inv.Wells = append(inv.Wells, *w)
			}
		}
	}

	if len(quality_ids) > 0 {
		q, err := assets.ReadString("pg/quality/by_inventory_quality_ids.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, quality_ids)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		issues := make([]model.Qualities, 0)
		rowsToStruct(rows, &issues)
		for _, iss := range issues {
			for _, inv := range q_inv[iss.ID] {
				inv.Qualities = append(inv.Qualities, iss)
			}
		}
	}

	return inventory, nil
}
