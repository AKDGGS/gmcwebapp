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

	inventory := model.Inventory{}
	c, err := ConnQuery(
		conn, "pg/inventory/by_inventory_id.sql",
		&inventory, id,
	)
	if err != nil {
		return nil, err
	}

	// If no inventory is found, stop right here
	if c == 0 {
		return nil, nil
	}

	if (flags & dbf.FILES) != 0 {
		_, err := ConnQuery(
			conn, "pg/file/by_inventory_id.sql",
			&inventory.Files, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		_, err = ConnQuery(
			conn, "pg/url/by_inventory_id.sql",
			&inventory.URLs, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		_, err = ConnQuery(
			conn, "pg/note/by_inventory_id.sql",
			&inventory.Notes, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.PUBLICATION) != 0 {
		_, err = ConnQuery(
			conn, "pg/publication/by_inventory_id.sql",
			&inventory.Publications, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.BOREHOLE) != 0 {
		_, err = ConnQuery(
			conn, "pg/borehole/by_inventory_id.sql",
			&inventory.Boreholes, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.OUTCROP) != 0 {
		_, err = ConnQuery(
			conn, "pg/outcrop/by_inventory_id.sql",
			&inventory.Outcrops, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.SHOTPOINT) != 0 {
		_, err = ConnQuery(
			conn, "pg/shotpoint/by_inventory_id.sql",
			&inventory.Shotpoints, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.WELL) != 0 {
		_, err = ConnQuery(
			conn, "pg/well/by_inventory_id.sql",
			&inventory.Wells, id,
		)
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
	}

	if (flags & dbf.QUALITY) != 0 {
		_, err = ConnQuery(
			conn, "pg/quality/by_inventory_id.sql",
			&inventory.Qualities, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.TRACKING) != 0 {
		_, err = ConnQuery(
			conn, "pg/container_log/by_inventory_id.sql",
			&inventory.ContainerLog, id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		q, err := assets.ReadString("pg/inventory/geojson.sql")
		if err != nil {
			return nil, err
		}

		row := conn.QueryRow(context.Background(), q, id)
		if err := row.Scan(&inventory.GeoJSON); err != nil && err != ErrNoRows {
			return nil, err
		}
	}
	return &inventory, nil
}
