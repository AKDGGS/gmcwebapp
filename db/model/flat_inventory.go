package model

import (
	"fmt"
)

type FlatInventory struct {
	ID             int32         `json:"id"`
	Collection     *string       `json:"collection,omitempty"`
	SampleNumber   *string       `json:"sample,omitempty"`
	SlideNumber    *string       `json:"slide,omitempty"`
	BoxNumber      *string       `json:"box,omitempty"`
	SetNumber      *string       `json:"set,omitempty"`
	CoreNumber     *string       `json:"core,omitempty"`
	CoreDiameter   *string       `json:"diameter,omitempty"`
	IntervalTop    *string       `json:"top,omitempty"`
	IntervalBottom *string       `json:"bottom,omitempty"`
	Keywords       *string       `json:"keywords,omitempty"`
	Barcode        *string       `json:"barcode,omitempty"`
	Remark         *string       `json:"remark,omitempty"`
	Geometries     []interface{} `json:"geometries,omitempty"`
	CanPublish     bool          `db:"can_publish" json:"can_publish"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
