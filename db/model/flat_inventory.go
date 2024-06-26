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
	Well           []struct {
		ID       int32   `json:"id"`
		Name     string  `json:"name"`
		AltNames *string `json:"altnames,omitempty"`
		Number   *string `json:"number,omitempty"`
		API      *string `json:"api,omitempty"`
	} `json:"well,omitempty"`
	Outcrop []struct {
		ID     int32   `json:"id"`
		Name   string  `json:"name"`
		Number *string `json:"number,omitempty"`
		Year   *int32  `json:"year,omitempty"`
	} `json:"outcrop,omitempty"`
	Borehole []struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		Prospect struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
			ARDF string `json:"ardf"`
		} `json:"prospect,omitempty"`
	} `json:"borehole,omitempty"`
	Shotlines []struct {
		ID   int32    `json:"id"`
		Name string   `json:"name"`
		Year *int32   `json:"year,omitempty"`
		Min  *float64 `json:"min,omitempty"`
		Max  *float64 `json:"max,omitempty"`
	} `json:"shotlines,omitempty"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
