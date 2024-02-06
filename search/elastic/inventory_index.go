package elastic

import (
	"gmc/db/model"
)

type ElasticInventoryIndex struct {
}

func (ii *ElasticInventoryIndex) Add(*model.FlatInventory) error {
	return nil
}

func (ii *ElasticInventoryIndex) Rollback() error {
	return nil
}

func (ii *ElasticInventoryIndex) Commit() error {
	return nil
}
