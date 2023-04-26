package pg

import (
	"context"

	"gmc/assets"
	"gmc/db/model"
)

func (pg *Postgres) GetFile(id int) (*model.File, error) {
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

func (pg *Postgres) PutFile(file *model.File, precommitFunc func() error) error {
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

	insert_wellfile_sql, err := assets.ReadString("pg/well/insert_wellfile.sql")
	if err != nil {
		return err
	}

	for _, well := range file.WellIDs {
		if well != 0 {
			_, err = tx.Exec(context.Background(), insert_wellfile_sql, fileID, well)
			if err != nil {
				return err
			}
		}
	}

	if err = precommitFunc(); err != nil {
		return err
	}

	// All files successfully added to the file table
	tx.Commit(context.Background())
	return nil
}
