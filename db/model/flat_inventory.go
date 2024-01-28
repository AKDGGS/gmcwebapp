package model

import (
	"encoding/json"
)

type FlatInventory struct {
	ID         int32         `json:"id"`
	Barcode    *string       `json:"barcode,omitempty"`
	Geometries []interface{} `json:"geometries,omitempty"`
}

func (f *FlatInventory) MarshalJSON() ([]byte, error) {
	type Alias FlatInventory
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}
