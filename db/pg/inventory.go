package pg

import (
	"context"

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

	if rowToStruct(rows, &inventory) == 0 {
		return nil, nil
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
		rowToStruct(rows, &inventory.Files)
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
		rowToStruct(r, &inventory.URLs)
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
		rowToStruct(rows, &inventory.Notes)
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
		rowToStruct(rows, &inventory.Publications)
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
		rowToStruct(rows, &inventory.Boreholes)
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
		rowToStruct(rows, &inventory.Outcrops)
	}

	if (flags & dbf.SHOTLINE) != 0 {
		q, err := assets.ReadString("pg/shotline/by_inventory_id.sql")
		if err != nil {
			return nil, err
		}
		rows, err := pg.pool.Query(context.Background(), q, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		rowToStruct(rows, &inventory.Shotpoints)
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
		rowToStruct(rows, &inventory.Wells)
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
		rowToStruct(rows, &inventory.Qualities)
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
		rowToStruct(rows, &inventory.ContainerLog)
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

func (pg *Postgres) GetInventoryByBarcode(barcode string, flags int) ([]model.Inventory, error) {
	q, err := assets.ReadString("pg/inventory/by_barcode.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inventory := []model.Inventory{}

	if rowToStruct(rows, &inventory) == 0 {
		return nil, nil
	}

	for i := 0; i < len(inventory); i++ {
		if (flags & dbf.BOREHOLE) != 0 {
			q, err := assets.ReadString("pg/borehole/by_inventory_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowToStruct(rows, &inventory[i].Boreholes)
		}

		if (flags & dbf.OUTCROP) != 0 {
			q, err := assets.ReadString("pg/outcrop/by_inventory_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowToStruct(rows, &inventory[i].Outcrops)
		}

		if (flags & dbf.SHOTPOINT) != 0 {
			q, err := assets.ReadString("pg/shotline/by_inventory_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowToStruct(rows, &inventory[i].Shotpoints)
		}

		if (flags & dbf.WELL) != 0 {
			q, err := assets.ReadString("pg/well/by_inventory_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowToStruct(rows, &inventory[i].Wells)
		}
		if (flags & dbf.QUALITY) != 0 {
			q, err := assets.ReadString("pg/quality/by_inventory_id.sql")
			if err != nil {
				return nil, err
			}
			rows, err := pg.pool.Query(context.Background(), q, inventory[i].ID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			rowToStruct(rows, &inventory[i].Qualities)
		}
	}
	return inventory, nil
}
