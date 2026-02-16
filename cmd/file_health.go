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
		fmt.Printf("Usage: %s %s [args]\n", exec, cmd)
		flagset.PrintDefaults()
	}
	fix_sz := flagset.Bool(
		"fix-size", false, "Fix invalid sizes in the database",
	)
	fix_hash := flagset.Bool(
		"fix-hash", false, "Fix invalid hashes in the database",
	)
	if err := flagset.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
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
		fp := fmt.Sprintf("%d/%s", f.ID, f.Name)
		stat, serr := fs.StatFile(fp)
		if serr != nil {
			failed = true
			if errors.Is(serr, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "file not found: %s\n", fp)
			}
		} else {
			dirty := false
			if *fix_sz && stat.Size != f.Size {
				f.Size = stat.Size
				dirty = true
			}

			if *fix_hash && stat.Hash != "" {
				if f.Hash == nil || stat.Hash != *f.Hash {
					f.Hash = &stat.Hash
					dirty = true
				}
			}

			if dirty {
				if derr := db.PutFile(f, nil); derr != nil {
					fmt.Printf("error updating file: %s\n", derr)
				}
			}
		}

		return nil
	})
	if err != nil || failed {
		return 1
	}
	return 0
}
