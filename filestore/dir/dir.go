package dir

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"mime"
	"os"
	"path/filepath"

	fsutil "gmc/filestore/util"
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

	md_b := big.NewInt(stat.ModTime().UnixMicro()).Bytes()
	sz_b := big.NewInt(stat.Size()).Bytes()

	return &fsutil.File{
		Name:         stat.Name(),
		Size:         stat.Size(),
		LastModified: stat.ModTime(),
		ContentType:  mt,
		Content:      file,
		ETag: fmt.Sprintf("%s-%s",
			base64.RawStdEncoding.EncodeToString(md_b),
			base64.RawStdEncoding.EncodeToString(sz_b),
		),
	}, nil
}

func (d *Dir) PutFile(f *fsutil.File) error {
	dir := filepath.Dir(filepath.Join(d.path, f.Name))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(d.path, f.Name))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, f.Content)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dir) Shutdown() {
	// Do nothing
}

func (d *Dir) DeleteFile(f *fsutil.File) error {
	err := os.Remove(filepath.Join(d.path, f.Name))
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	err = os.Remove(filepath.Dir(filepath.Join(d.path, f.Name)))
	if err != nil {
		return fmt.Errorf("error deleting directory: %w", err)
	}
	return nil
}
