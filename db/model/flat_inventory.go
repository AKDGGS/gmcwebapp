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

func FlatInventoryFields(full bool) []string {
	if full {
		return []string{
			"ID",
			"Related",
			"Project ID",
			"Project",
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
			"Container",
			"Description",
		}
	}
	return []string{
		"ID",
		"Related",
		"Project ID",
		"Project",
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
		"Description",
	}
}

func (f *FlatInventory) AsStringArray(full bool) []string {
	var rel strings.Builder
	for _, w := range f.Well {
		if rel.Len() > 0 {
			rel.WriteString("\n")
		}
		rel.WriteString("Well: ")
		rel.WriteString(w.DisplayName)
		if w.Number != nil {
			rel.WriteString(" - ")
			rel.WriteString(*w.Number)
		}
		if w.API != nil {
			rel.WriteString("\nAPI: ")
			rel.WriteString(*w.API)
		}
	}
	for _, o := range f.Outcrop {
		if rel.Len() > 0 {
			rel.WriteString("\n")
		}
		rel.WriteString("Outcrop: ")
		rel.WriteString(o.Name)
		if o.Number != nil {
			rel.WriteString(" - ")
			rel.WriteString(*o.Number)
		}
	}
	for _, b := range f.Borehole {
		if rel.Len() > 0 {
			rel.WriteString("\n")
		}
		if b.Prospect.DisplayName != "" {
			rel.WriteString("Prospect: ")
			rel.WriteString(b.Prospect.DisplayName)
			rel.WriteString("\n")
		}
		rel.WriteString("Borehole: ")
		rel.WriteString(b.DisplayName)
	}
	for _, s := range f.Shotline {
		if rel.Len() > 0 {
			rel.WriteString("\n")
		}
		rel.WriteString("Shotline: ")
		rel.WriteString(s.Name)
		if s.Max != nil {
			rel.WriteString("\nShotpoints: ")
			rel.WriteString(strconv.FormatFloat(*s.Min, 'f', 2, 64))
			rel.WriteString(" - ")
			rel.WriteString(strconv.FormatFloat(*s.Max, 'f', 2, 64))
		}
	}
	for _, p := range f.Publication {
		if rel.Len() > 0 {
			rel.WriteString("\n")
		}
		rel.WriteString("Publication: ")
		rel.WriteString(p.Title)
		if p.Year != nil {
			rel.WriteString(" (")
			rel.WriteString(strconv.FormatInt(int64(*p.Year), 10))
			rel.WriteString(")")
		}
		if p.Description != "" {
			rel.WriteString("\nDescription: ")
			rel.WriteString(p.Description)
		}
	}

	if full {
		return []string{
			qfmt(f.ID),
			rel.String(),
			qfmt(f.ProjectID),
			qfmt(f.Project),
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
			qfmt(f.ContainerPath),
			qfmt(f.Description),
		}
	}
	return []string{
		qfmt(f.ID),
		rel.String(),
		qfmt(f.ProjectID),
		qfmt(f.Project),
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
		qfmt(f.Description),
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
