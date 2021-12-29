package filestore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Dir struct {
	path string
}

func newDir(cfg map[string]interface{}) (*Dir, error) {
	path, ok := cfg["path"].(string)
	if !ok {
		return nil, fmt.Errorf("directory path must exist and be a string")
	}

	return &Dir{path: path}, nil
}

func (d *Dir) GetFile(name string) (io.ReadSeekCloser, error) {
	f, err := os.Open(filepath.Join(d.path, name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *Dir) Shutdown() {
	// Do nothing
}
