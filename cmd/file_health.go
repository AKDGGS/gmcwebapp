package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

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
	failed := false
	err = db.ListFiles(func(f *model.File) error {
		fp := fmt.Sprintf("%d/%s", f.ID, *f.Name)
		if _, serr := fs.StatFile(fp); serr != nil {
			failed = true
			if errors.Is(serr, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "file not found: %s\n", fp)
			}
		}
		return nil
	})
	if err != nil || failed {
		return 1
	}
	return 0
}
