package filestore

import (
	"fmt"
	"gmc/config"
	"io"
	"time"
)

type FileStore interface {
	GetFile(string) (*File, error)
	Shutdown()
}

type File struct {
	Name         string
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
	Content      io.ReadSeekCloser
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
		return nil, fmt.Errorf("file_store type cannot be empty")
	default:
		return nil, fmt.Errorf("unknown file store type: %s", cfg.Type)
	}
	return stor, nil
}
