package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"gmc/config"
	"gmc/search"
	searchutil "gmc/search/util"
)

func SearchCommand(cfg *config.Config, exec, cmd string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "query argument required\n")
		return 1
	}

	sea, err := search.New(cfg.Search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	params := &searchutil.InventoryParams{
		Size:  10,
		Query: args[0],
	}

	r, err := sea.SearchInventory(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	jsn, err := json.Marshal(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	fmt.Println(string(jsn))
	return 0
}
