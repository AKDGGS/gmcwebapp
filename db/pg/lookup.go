package pg

import (
	"context"

	"gmc/assets"
)

func (pg *Postgres) GetBarcodeCounts(barcode string) (int, int, error) {
	q, err := assets.ReadString("pg/lookup/barcode_counts.sql")
	if err != nil {
		return 0, 0, err
	}
	var container_count, inventory_count int
	err = pg.pool.QueryRow(context.Background(), q, barcode).Scan(&inventory_count, &container_count)
	if err != nil {
		return 0, 0, err
	}
	return inventory_count, container_count, nil
}
