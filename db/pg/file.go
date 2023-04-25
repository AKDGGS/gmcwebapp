package pg

import (
	"context"
	"fmt"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetFile(id int) (*model.File, error) {
	fmt.Println("GetFile")
	q, err := assets.ReadString("pg/file/by_file_id.sql")
	if err != nil {
		return nil, err
	}
	rows, err := pg.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	file := model.File{}
	if rowToStruct(rows, &file) == 0 {
		return nil, nil
	}
	return &file, nil
}

func (pg *Postgres) PutFile(file *model.File, callback func() error) error {
	tx, err := pg.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	insert_sql, err := assets.ReadString("pg/file/insert.sql")
	if err != nil {
		return err
	}

	var fileID int

	err = tx.QueryRow(context.Background(), insert_sql, file.Name, file.Description,
		file.Size, file.Type, file.MD5).Scan(&fileID)
	if err != nil {
		return err
	}
	file.ID = fileID

	if err = callback(); err != nil {
		return err
	}

	// All files successfully added to the file table
	tx.Commit(context.Background())
	return nil
}
