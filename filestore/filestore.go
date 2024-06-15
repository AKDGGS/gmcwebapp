package filestore

import (
	"fmt"

	"gmc/config"
	"gmc/filestore/dir"
	"gmc/filestore/s3"
	fsutil "gmc/filestore/util"
)

type FileStore interface {
	GetFile(string) (*fsutil.File, error)
	PutFile(*fsutil.File) error
	DeleteFile(*fsutil.File) error
	Shutdown()
}

func New(cfg config.FileStoreConfig) (FileStore, error) {
	var stor FileStore
	var err error
	switch cfg.Type {
	case "s3", "minio":
		stor, err = s3.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "dir", "directory":
		stor, err = dir.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown file store type: %s", cfg.Type)
	}
	return stor, nil
}
