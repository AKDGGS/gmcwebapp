package cmd

import (
	"fmt"
	"os"
	"strconv"

	"gmc/config"
	"gmc/db"
	"gmc/filestore"
)

func FileDelete(exec string, cfg *config.Config, cmd string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%s: file removal requires an ID\n", exec)
		return 1
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: \"%s\" is in a valid ID\n", exec, args[9])
		return 1
	}

	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	fs, err := filestore.NewFileStores(cfg.FileStores)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	file, err := db.GetFile(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	if err := db.DeleteFile(id, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	err = fs.DeleteFile(fmt.Sprintf("%d/%s", file.ID, file.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	return 0
}
