package pg

import (
	"context"

	dbf "gmc/db/flag"
	"gmc/db/model"
)

func (pg *Postgres) GetInventory(id int, flags int) (*model.Inventory, error) {
	conn, err := pg.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	inventory, err := cQryStruct[model.Inventory](
		conn, "pg/inventory/by_inventory_id.sql", id,
	)
	if err != nil {
		return nil, err
	}

	// If no inventory is found, stop right here
	if inventory == nil {
		return nil, nil
	}

	if (flags & dbf.WELL) != 0 {
		inventory.Wells, err = cQryStructs[model.Well](
			conn, "pg/well/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.NOTE) != 0 {
		inventory.Notes, err = cQryStructs[model.Note](
			conn, "pg/note/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.URLS) != 0 {
		inventory.URLs, err = cQryStructs[model.URL](
			conn, "pg/url/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.QUALITY) != 0 {
		inventory.Qualities, err = cQryStructs[model.Quality](
			conn, "pg/quality/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.FILES) != 0 {
		inventory.Files, err = cQryStructs[model.File](
			conn, "pg/file/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.PUBLICATION) != 0 {
		inventory.Publications, err = cQryStructs[model.Publication](
			conn, "pg/publication/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.BOREHOLE) != 0 {
		inventory.Boreholes, err = cQryStructs[model.Borehole](
			conn, "pg/borehole/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.OUTCROP) != 0 {
		inventory.Outcrops, err = cQryStructs[model.Outcrop](
			conn, "pg/outcrop/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.TRACKING) != 0 {
		inventory.ContainerLog, err = cQryStructs[model.ContainerLog](
			conn, "pg/container_log/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.SHOTPOINT) != 0 {
		inventory.Shotpoints, err = cQryStructs[model.Shotpoint](
			conn, "pg/shotpoint/by_inventory_id.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}

	if (flags & dbf.GEOJSON) != 0 {
		inventory.GeoJSON, err = cQryValue(
			conn, "pg/inventory/geojson.sql", id,
		)
		if err != nil {
			return nil, err
		}
	}
	return inventory, nil
}
