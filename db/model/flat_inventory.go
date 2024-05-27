package model

import (
	"fmt"
)

type FlatInventory struct {
	ID         int32         `json:"id"`
	Collection *string       `json:"collection,omitempty"`
	Barcode    *string       `json:"barcode,omitempty"`
	Remark     *string       `json:"remark,omitempty"`
	Geometries []interface{} `json:"geometries,omitempty"`
	CanPublish bool          `db:"can_publish" json:"can_publish"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
