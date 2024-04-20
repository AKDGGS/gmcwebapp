package cmd

import (
	"fmt"
	"os"
	"time"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
	"gmc/search"
)

func IndexCommand(cfg *config.Config, exec, cmd string, args []string) int {
	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		return 1
	}

	sea, err := search.New(cfg.Search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	ii, err := sea.NewInventoryIndex()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err.Error())
		return 1
	}

	fmt.Printf("indexing ")
	start := time.Now()
	err = db.GetFlatInventory(func(f *model.FlatInventory) error {
		if ii.Count() > 1000 {
			if err := ii.Flush(); err != nil {
				return err
			}
			fmt.Print(".")
		}
		return ii.Add(f)
	})
	if err != nil {
		ii.Rollback()
		fmt.Fprintf(os.Stderr, "\n%s: %s\n",
			os.Args[0], err.Error())
		return 1
	}

	if err := ii.Commit(); err != nil {
		ii.Rollback()
		fmt.Fprintf(os.Stderr, "\n%s: %s\n", os.Args[0], err.Error())
		return 1
	}
	fmt.Printf(" complete in %s\n", time.Since(start).String())
	return 0
}
