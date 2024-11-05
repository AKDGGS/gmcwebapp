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
	IntervalTop    *float64        `db:"interval_top" json:"top,omitempty"`
	IntervalBottom *float64        `db:"interval_bottom" json:"bottom,omitempty"`
	IntervalUnit   *string         `db:"interval_unit" json:"unit,omitempty"`
	Keyword        []string        `json:"keyword,omitempty"`
	Barcode        []string        `json:"barcode,omitempty"`
	DisplayBarcode *string         `json:"display_barcode,omitempty"`
	ContainerID    *int32          `json:"container_id,omitempty"`
	ContainerPath  *string         `json:"path_cache,omitempty"`
	Remark         *string         `json:"remark,omitempty"`
	Geometries     json.RawMessage `json:"geometries,omitempty"`
	ProjectID      *int32          `json:"project_id,omitempty"`
	Project        *string         `json:"project,omitempty"`
	CanPublish     *bool           `db:"can_publish" json:"can_publish,omitempty"`
	Description    *string         `json:"description,omitempty"`
	Note           []string        `json:"note,omitempty"`
	Issue          []string        `json:"issue,omitempty"`
	Interval       *struct {
		GTE *float64 `db:"gte" json:"gte,omitempty"`
		LTE *float64 `db:"lte" json:"lte,omitempty"`
	} `json:"interval,omitempty"`
	Well []struct {
		ID          int32    `json:"id"`
		DisplayName string   `json:"display_name"`
		Name        []string `json:"name,omitempty"`
		Number      *string  `json:"number,omitempty"`
		API         *string  `json:"api,omitempty"`
	} `json:"well,omitempty"`
	Outcrop []struct {
		ID     int32   `json:"id"`
		Name   string  `json:"name"`
		Number *string `json:"number,omitempty"`
		Year   *int32  `json:"year,omitempty"`
	} `json:"outcrop,omitempty"`
	Borehole []struct {
		ID          int32    `json:"id"`
		DisplayName string   `json:"display_name"`
		Name        []string `json:"name,omitempty"`
		Prospect    struct {
			ID          int32    `json:"id"`
			DisplayName string   `json:"display_name"`
			Name        []string `json:"name,omitempty"`
			ARDF        string   `json:"ardf,omitempty"`
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
		Year        *int32 `json:"year,omitempty"`
		Description string `json:"description,omitempty"`
		Number      string `json:"number,omitempty"`
		Series      string `json:"series,omitempty"`
	} `json:"publication,omitempty"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}
