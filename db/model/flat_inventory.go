package model

import (
	"encoding/json"
	"fmt"
)

type FlatInventory struct {
	ID             int32           `json:"id"`
	Collection     *string         `json:"collection,omitempty"`
	CollectionID   *int32          `json:"collection_id,omitempty"`
	SampleNumber   *string         `json:"sample,omitempty"`
	SlideNumber    *string         `json:"slide,omitempty"`
	BoxNumber      *string         `json:"box,omitempty"`
	SetNumber      *string         `json:"set,omitempty"`
	CoreNumber     *string         `json:"core,omitempty"`
	CoreDiameter   *float64        `json:"diameter,omitempty"`
	CoreName       *string         `json:"core_name,omitempty"`
	CoreUnit       *string         `json:"core_unit,omitempty"`
	IntervalTop    *float64        `json:"top,omitempty"`
	IntervalBottom *float64        `json:"bottom,omitempty"`
	IntervalUnit   *string         `json:"interval_unit,omitempty"`
	Keywords       []string        `json:"keywords,omitempty"`
	Barcode        []string        `json:"barcode,omitempty"`
	DisplayBarcode *string         `json:"display_barcode,omitempty"`
	ContainerID    *int32          `json:"container_id,omitempty"`
	ContainerPath  *string         `json:"path_cache,omitempty"`
	Remark         *string         `json:"remark,omitempty"`
	Geometries     json.RawMessage `json:"geometries,omitempty"`
	ProjectID      *int32          `json:"project_id,omitempty"`
	Project        *string         `json:"project,omitempty"`
	CanPublish     bool            `db:"can_publish" json:"can_publish"`
	Description    *string         `json:"description,omitempty"`
	Note           []string        `json:"note,omitempty"`
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
	Shotline []struct {
		ID   int32    `json:"id"`
		Name string   `json:"name"`
		Year *int32   `json:"year,omitempty"`
		Min  *float64 `json:"min,omitempty"`
		Max  *float64 `json:"max,omitempty"`
	} `json:"shotline,omitempty"`
	Publication []struct {
		ID          int32  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Number      string `json:"number,omitempty"`
		Series      string `json:"series,omitempty"`
	} `json:"publication,omitempty"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
