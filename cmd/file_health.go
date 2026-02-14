package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/filestore"
)

func FileHealth(exec string, cfg *config.Config, cmd string, args []string) int {
	flagset := flag.NewFlagSet(cmd, flag.ContinueOnError)
	flagset.SetOutput(os.Stdout)
	flagset.Usage = func() {
		fmt.Printf("Usage: %s %s\n",
			exec, cmd)
		flagset.PrintDefaults()
	}
	fs, err := filestore.NewFileStores(cfg.FileStores)
	if err != nil {
		return 1
	}

	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		return 1
	}
	err = db.ListFiles(func(f *model.File) error {
		if f.Name != nil {
			id := strconv.Itoa(int(f.ID))
			_, innererr := fs.StatFile(id + "/" + *f.Name)
			if innererr != nil {
				for _, curerr := range innererr.(interface{ Unwrap() []error }).Unwrap() {
					if !errors.Is(curerr, os.ErrNotExist) {
						return innererr
					}
				}
				fmt.Fprintf(os.Stderr, "%s: file not found (%d): %s\n", exec, f.ID, *f.Name)
			}
		}
		return nil
	})
	if err != nil {
		return 1
	}
	return 0
}
