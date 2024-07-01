package pg

import (
	"context"
	"fmt"
	"strings"

	"gmc/assets"
)

func (pg *Postgres) AddAudit(remark string, container_list []string) error {
	if remark != "" {
		remark = strings.TrimSpace(remark)
	}
	if remark == "" && len(container_list) == 0 {
		return fmt.Errorf("remark and the items list cannot be empty")
	}
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q, err := assets.ReadString("pg/audit/insert_group.sql")
	if err != nil {
		return err
	}
	var audit_group_id int32
	err = tx.QueryRow(context.Background(), q, remark).Scan(&audit_group_id)
	if err != nil {
		if err == ErrNoRows {
			return nil
		}
		return err
	}
	if audit_group_id == 0 {
		return fmt.Errorf("audit insert returned zero for audit_group_id")
	}
	for _, c := range container_list {
		q, err = assets.ReadString("pg/audit/insert.sql")
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), q, audit_group_id, c)
		if err != nil {
			return err
		}
	}
	// If the insert is successful, commit the changes
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}
