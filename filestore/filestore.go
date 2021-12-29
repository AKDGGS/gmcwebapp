package filestore

import (
	"fmt"
	"gmc/config"
	"io"
)

type FileStore interface {
	GetFile(string) (io.ReadSeekCloser, error)
	Shutdown()
}

func New(cfg config.FileStoreConfig) (FileStore, error) {
	var stor FileStore
	var err error
	switch cfg.Type {
	case "s3", "minio":
		stor, err = newS3(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "dir", "directory":
		stor, err = newDir(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "":
		return nil, fmt.Errorf("file_store type may not be empty")
	default:
		return nil, fmt.Errorf("Unknown file store type: %s", cfg.Type)
	}
	return stor, nil
}
