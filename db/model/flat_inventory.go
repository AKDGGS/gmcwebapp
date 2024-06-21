package model

import (
	"encoding/json"
	"fmt"
)

type FlatInventory struct {
	ID             int32           `json:"id"`
	Collection     *string         `json:"collection,omitempty"`
	SampleNumber   *string         `json:"sample,omitempty"`
	SlideNumber    *string         `json:"slide,omitempty"`
	BoxNumber      *string         `json:"box,omitempty"`
	SetNumber      *string         `json:"set,omitempty"`
	CoreNumber     *string         `json:"core,omitempty"`
	CoreDiameter   *float64        `json:"diameter,omitempty"`
	IntervalTop    *float64        `json:"top,omitempty"`
	IntervalBottom *float64        `json:"bottom,omitempty"`
	Keywords       []string        `json:"keywords,omitempty"`
	Barcode        *string         `json:"barcode,omitempty"`
	Remark         *string         `json:"remark,omitempty"`
	Geometries     json.RawMessage `json:"geometries,omitempty"`
	CanPublish     bool            `db:"can_publish" json:"can_publish"`
	Wells          []struct {
		ID       int32   `json:"id"`
		Name     string  `json:"name"`
		AltNames *string `json:"altnames,omitempty"`
		Number   *string `json:"number,omitempty"`
		API      *string `json:"api,omitempty"`
	} `json:"wells,omitempty"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
