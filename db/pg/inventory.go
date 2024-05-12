package pg

import (
	"context"

	"gmc/assets"
	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetInventory(id int, flags int) (*model.Inventory, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	q, err := assets.ReadString("pg/inventory/by_inventory_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := model.Inventory{}
	c, err := rowsToStruct(rows, &inventory)
	if err != nil {
		return nil, err
	}
	rows.Close()

	//nothing returned by the database
	if c == 0 {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		q, err = assets.ReadString("pg/file/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Files)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.URLS) != 0 {
		q, err = assets.ReadString("pg/url/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		r, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(r, &inventory.URLs)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.NOTE) != 0 {
		q, err = assets.ReadString("pg/note/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Notes)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.PUBLICATION) != 0 {
		q, err = assets.ReadString("pg/publication/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Publications)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.BOREHOLE) != 0 {
		q, err := assets.ReadString("pg/borehole/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Boreholes)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.OUTCROP) != 0 {
		q, err := assets.ReadString("pg/outcrop/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Outcrops)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.SHOTPOINT) != 0 {
		q, err := assets.ReadString("pg/shotpoint/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Shotpoints)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.WELL) != 0 {
		q, err := assets.ReadString("pg/well/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Wells)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.ORGANIZATION) != 0 {
		for i := 0; i < len(inventory.Boreholes); i++ {
			q, err := assets.ReadString("pg/organization/by_borehole_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := conn.Query(context.Background(), q, inventory.Boreholes[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			_, err = rowsToStruct(rows, &inventory.Boreholes[i].Organizations)
			if err != nil {
				return nil, err
			}
			rows.Close()
		}
		for i := 0; i < len(inventory.Outcrops); i++ {
			q, err := assets.ReadString("pg/organization/by_outcrop_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := conn.Query(context.Background(), q, inventory.Outcrops[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			_, err = rowsToStruct(rows, &inventory.Outcrops[i].Organizations)
			if err != nil {
				return nil, err
			}
			rows.Close()
		}
		for i := 0; i < len(inventory.Wells); i++ {
			q, err := assets.ReadString("pg/organization/by_well_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := conn.Query(context.Background(), q, inventory.Wells[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			_, err = rowsToStruct(rows, &inventory.Wells[i].Organizations)
			if err != nil {
				return nil, err
			}
			rows.Close()
		}
	}

	if (flags & dbf.QUALITY) != 0 {
		q, err := assets.ReadString("pg/quality/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.Qualities)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.TRACKING) != 0 {
		q, err := assets.ReadString("pg/container_log/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := conn.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		_, err = rowsToStruct(rows, &inventory.ContainerLog)
		if err != nil {
			return nil, err
		}
		rows.Close()
	}

	if (flags & dbf.GEOJSON) != 0 {
		q, err := assets.ReadString("pg/inventory/geojson.sql")
		if err != nil {
			return nil, err
		}

		row := conn.QueryRow(context.Background(), q, id)
		if err := row.Scan(&inventory.GeoJSON); err != nil {
			return nil, err
		}
	}
	return &inventory, nil
}
