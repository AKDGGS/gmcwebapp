package search

import (
	"fmt"

	"gmc/config"
	"gmc/search/elastic"
	"gmc/search/util"
)

type Search interface {
	Name() string
	NewInventoryIndex() (util.InventoryIndex, error)
	InventorySortByFields() [][2]string
	SearchInventory(*util.InventoryParams) (*util.InventoryResults, error)
	Shutdown()
}

func New(cfg config.SearchConfig) (Search, error) {
	var sear Search
	var err error
	switch cfg.Type {
	case "elastic", "elasticsearch":
		sear, err = elastic.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "":
		return nil, fmt.Errorf("search type may not be empty")
	default:
		return nil, fmt.Errorf("unknown search type: %s", cfg.Type)
	}
	return sear, nil
}
