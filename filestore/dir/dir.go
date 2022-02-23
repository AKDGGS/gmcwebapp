package dir

import (
	"fmt"
	fsutil "gmc/filestore/util"
	"mime"
	"os"
	"path/filepath"
)

type Dir struct {
	path string
}

func New(cfg map[string]interface{}) (*Dir, error) {
	path, ok := cfg["path"].(string)
	if !ok {
		return nil, fmt.Errorf("directory path must exist and be a string")
	}

	return &Dir{path: path}, nil
}

func (d *Dir) GetFile(name string) (*fsutil.File, error) {
	fp := filepath.Join(d.path, name)
	if !filepath.IsAbs(fp) {
		return nil, os.ErrPermission
	}

	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	mt := mime.TypeByExtension(filepath.Ext(stat.Name()))
	if mt == "" {
		mt = "application/octet-stream"
	}

	return &fsutil.File{
		Name:         stat.Name(),
		ETag:         "",
		Size:         stat.Size(),
		LastModified: stat.ModTime(),
		ContentType:  mt,
		Content:      file,
	}, nil
}

func (d *Dir) Shutdown() {
	// Do nothing
}
