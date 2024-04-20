package search

import (
	"fmt"

	"gmc/config"
	"gmc/search/elastic"
	"gmc/search/util"
)

type Search interface {
	NewInventoryIndex() (util.InventoryIndex, error)
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
		return nil, fmt.Errorf("search_provider type may not be empty")
	default:
		return nil, fmt.Errorf("Unknown search provider type: %s", cfg.Type)
	}
	return sear, nil
}
