package cmd

import (
	"fmt"
	"os"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
)

func IndexCommand(cfg *config.Config, exec, cmd string, args []string) int {
	db, err := db.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		return 1
	}

	err = db.GetFlatInventory(func(f *model.FlatInventory) error {
		js, err := f.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(js))
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		return 1
	}
	return 0
}
