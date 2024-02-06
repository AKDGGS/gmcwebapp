package util

import (
	"gmc/db/model"
)

type InventoryIndex interface {
	Add(*model.FlatInventory) error
	Rollback() error
	Commit() error
}
