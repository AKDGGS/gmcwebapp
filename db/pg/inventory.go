package pg

import (
	"context"

	"gmc/assets"
	dbe "gmc/db/errors"
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

	c, err := rowsToStruct(rows, &inventory)
	if err != nil {
		return nil, err
	}

	//nothing returned by the database
	if c == 0 {
		return nil, dbe.ErrBarcodeNotFound
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
		_, err = rowsToStruct(rows, &inventory.Files)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(r, &inventory.URLs)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Notes)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Publications)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Boreholes)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Outcrops)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Shotpoints)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.Wells)
		if err != nil {
			return nil, err
		}
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
			_, err = rowsToStruct(rows, &inventory.Boreholes[i].Organizations)
			if err != nil {
				return nil, err
			}
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
			_, err = rowsToStruct(rows, &inventory.Outcrops[i].Organizations)
			if err != nil {
				return nil, err
			}
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
			_, err = rowsToStruct(rows, &inventory.Wells[i].Organizations)
			if err != nil {
				return nil, err
			}
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
		_, err = rowsToStruct(rows, &inventory.Qualities)
		if err != nil {
			return nil, err
		}
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
		_, err = rowsToStruct(rows, &inventory.ContainerLog)
		if err != nil {
			return nil, err
		}
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
