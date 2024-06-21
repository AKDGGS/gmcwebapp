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
	Query   string
	From    int
	Size    int
	Private bool
}

type InventoryResults struct {
	Hits  []InventoryHit `json:"hits,omitempty"`
	From  int            `json:"from"`
	Total int64          `json:"total"`
	Time  time.Duration  `json:"time"`
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

type InventoryHit struct {
	ID             int             `json:"id"`
	Collection     string          `json:"collection,omitempty"`
	SampleNumber   string          `json:"sample,omitempty"`
	SlideNumber    string          `json:"slide,omitempty"`
	BoxNumber      string          `json:"box,omitempty"`
	SetNumber      string          `json:"set,omitempty"`
	CoreNumber     string          `json:"core,omitempty"`
	CoreDiameter   *float64        `json:"diameter,omitempty"`
	IntervalTop    *float64        `json:"top,omitempty"`
	IntervalBottom *float64        `json:"bottom,omitempty"`
	Keywords       []string        `json:"keywords,omitempty"`
	Barcode        string          `json:"barcode,omitempty"`
	CanPublish     bool            `json:"can_publish"`
	Geometries     json.RawMessage `json:"geometries,omitempty"`
	Wells          []struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		AltNames *string `json:"altnames,omitempty"`
		Number   *string `json:"number,omitempty"`
		API      *string `json:"api,omitempty"`
	} `json:"wells,omitempty"`
}
