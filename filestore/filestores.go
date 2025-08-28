package filestore

import (
	"errors"
	"fmt"
	"strings"

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
	var errlist strings.Builder
	for _, store := range fs.stores {
		file, err := store.GetFile(name)
		if err == nil {
			return file, nil
		}
		errlist.WriteString("\n")
		errlist.WriteString(err.Error())
	}
	return nil, errors.New(errlist.String())
}

func (fs FileStores) PutFile(file *fsutil.File) error {
	if len(fs.stores) == 0 {
		return fmt.Errorf("no filestores configured")
	}
	var errlist strings.Builder
	for _, store := range fs.stores {
		if err := store.PutFile(file); err != nil {
			errlist.WriteString("\n")
			errlist.WriteString(err.Error())
		} else {
			return nil
		}
	}
	return errors.New(errlist.String())
}

func (fs FileStores) DeleteFile(file *fsutil.File) error {
	if len(fs.stores) == 0 {
		return fmt.Errorf("no filestores configured")
	}
	var errlist strings.Builder
	for _, store := range fs.stores {
		if err := store.DeleteFile(file); err != nil {
			errlist.WriteString("\n")
			errlist.WriteString(err.Error())
		} else {
			return nil
		}
	}
	return errors.New(errlist.String())
}

func (fs FileStores) Shutdown() {
	for _, store := range fs.stores {
		store.Shutdown()
	}
}
