package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func FlatInventoryFields() []string {
	return []string{
		"ID",
		"Collection ID",
		"Collection",
		"Sample Number",
		"Slide Number",
		"Box Number",
		"Set Number",
		"Core Number",
		"Core Diameter",
		"Core Name",
		"Core Unit",
		"Interval Top",
		"Interval Bottom",
		"Interval Unit",
		"Keywords",
		"Barcode",
		"Container ID",
		"Container",
		"Project ID",
		"Project",
	}
}

func (f *FlatInventory) AsStringArray() []string {
	return []string{
		qfmt(f.ID),
		qfmt(f.CollectionID),
		qfmt(f.Collection),
		qfmt(f.SampleNumber),
		qfmt(f.SlideNumber),
		qfmt(f.BoxNumber),
		qfmt(f.SetNumber),
		qfmt(f.CoreNumber),
		qfmt(f.CoreDiameter),
		qfmt(f.CoreName),
		qfmt(f.CoreUnit),
		qfmt(f.IntervalTop),
		qfmt(f.IntervalBottom),
		qfmt(f.IntervalUnit),
		qfmt(f.Keyword),
		qfmt(f.DisplayBarcode),
		qfmt(f.ContainerID),
		qfmt(f.ContainerPath),
		qfmt(f.ProjectID),
		qfmt(f.Project),
	}
}

func qfmt(v interface{}) string {
	switch t := v.(type) {
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case *int32:
		if t == nil {
			return ""
		}
		return strconv.FormatInt(int64(*t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case *int64:
		if t == nil {
			return ""
		}
		return strconv.FormatInt(*t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', 2, 64)
	case *float64:
		if t == nil {
			return ""
		}
		return strconv.FormatFloat(*t, 'f', 2, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', 2, 32)
	case *float32:
		if t == nil {
			return ""
		}
		return strconv.FormatFloat(float64(*t), 'f', 2, 32)
	case string:
		return t
	case *string:
		if t == nil {
			return ""
		}
		return *t
	case []string:
		return strings.Join(t, "; ")
	default:
		return ""
	}
}
