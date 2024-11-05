package util

import (
	"encoding/json"
	"time"

	"gmc/db/model"
)

type InventoryIndex interface {
	Add(*model.FlatInventory) error
	Count() int
	Commit() error
	Flush() error
	Rollback() error
}

type InventoryParams struct {
	Query          string
	Keywords       []string
	ProspectIDs    []int
	CollectionIDs  []int
	IntervalTop    *float64
	IntervalBottom *float64
	From           int
	Size           int
	Private        bool
	Sort           [][2]string
}

type InventoryResults struct {
	Hits    []model.FlatInventory `json:"hits,omitempty"`
	From    int                   `json:"from"`
	Total   int64                 `json:"total"`
	Time    time.Duration         `json:"time"`
	Private bool                  `json:"private,omitempty"`
}

func (ir *InventoryResults) MarshalJSON() ([]byte, error) {
	type Alias InventoryResults

	time := ir.Time.String()
	return json.Marshal(&struct {
		Time string `json:"time"`
		*Alias
	}{
		Time:  time,
		Alias: (*Alias)(ir),
	})
}
