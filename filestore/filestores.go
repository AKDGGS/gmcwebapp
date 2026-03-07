package filestore

import (
	"errors"
	"fmt"

	"gmc/config"
	fsutil "gmc/filestore/util"
)

type FileStores struct {
	stores []FileStore
}

func NewFileStores(configs *[]config.FileStoreConfig) (FileStores, error) {
	fs := FileStores{}
	if configs == nil {
		return fs, nil
	}
	for _, config := range *configs {
		store, err := NewStore(config)
		if err != nil {
			return fs, err
		}
		fs.stores = append(fs.stores, store)
	}
	return fs, nil
}

func (fs FileStores) GetFile(name string) (*fsutil.File, error) {
	if len(fs.stores) == 0 {
		return nil, fmt.Errorf("no filestores configured")
	}
	var errlist error
	for _, store := range fs.stores {
		file, err := store.GetFile(name)
		if err == nil {
			return file, nil
		}
		errlist = errors.Join(errlist, err)
	}
	return nil, errlist
}

func (fs FileStores) PutFile(file *fsutil.File) error {
	if len(fs.stores) == 0 {
		return fmt.Errorf("no filestores configured")
	}
	var errlist error
	for _, store := range fs.stores {
		if err := store.PutFile(file); err != nil {
			errlist = errors.Join(errlist, err)
		} else {
			return nil
		}
	}
	return errlist
}

func (fs FileStores) DeleteFile(name string) error {
	if len(fs.stores) == 0 {
		return fmt.Errorf("no filestores configured")
	}
	var errlist error
	for _, store := range fs.stores {
		if err := store.DeleteFile(name); err != nil {
			errlist = errors.Join(errlist, err)
		} else {
			return nil
		}
	}
	return errlist
}

func (fs FileStores) StatFile(name string) (*fsutil.FileStat, error) {
	if len(fs.stores) == 0 {
		return nil, fmt.Errorf("no filestores configured")
	}
	var errlist error
	for _, store := range fs.stores {
		stat, err := store.StatFile(name)
		if err != nil {
			errlist = errors.Join(errlist, err)
		} else {
			return stat, nil
		}
	}
	return nil, errlist
}

func (fs FileStores) Shutdown() {
	for _, store := range fs.stores {
		store.Shutdown()
	}
}
