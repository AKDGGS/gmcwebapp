package pg

import (
	"context"
	"fmt"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) AddContainer(barcode string, name string, remark string) error {
	if barcode == "" || len(strings.TrimSpace(barcode)) < 1 {
		return fmt.Errorf("Barcode cannot be empty")
	}
	q, err := assets.ReadString("pg/container/insert.sql")
	if err != nil {
		return err
	}
	_, err = pg.pool.Exec(context.Background(), q, barcode, name, remark)
	if err != nil {
		return err
	}
	return nil
}
