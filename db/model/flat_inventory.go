package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type FlatInventory struct {
	ID             int32             `json:"id" parquet:"id,uncompressed"`
	Collection     *string           `json:"collection,omitempty" parquet:"collection,uncompressed,optional"`
	CollectionID   *int32            `json:"collection_id,omitempty" parquet:"-"`
	SampleNumber   *string           `json:"sample,omitempty" parquet:"sample,uncompressed,optional"`
	SlideNumber    *string           `json:"slide,omitempty" parquet:"slide,uncompressed,optional"`
	BoxNumber      *string           `json:"box,omitempty" parquet:"box,uncompressed,optional"`
	SetNumber      *string           `json:"set,omitempty" parquet:"set,uncompressed,optional"`
	CoreNumber     *string           `json:"core,omitempty" parquet:"core,uncompressed,optional"`
	CoreDiameter   *float64          `json:"diameter,omitempty" parquet:"diameter,uncompressed,optional"`
	CoreName       *string           `json:"core_name,omitempty" parquet:"core_name,uncompressed,optional"`
	CoreUnit       *string           `json:"core_unit,omitempty" parquet:"core_unit,uncompressed,optional"`
	IntervalTop    *float64          `db:"interval_top" json:"top,omitempty" parquet:"top,uncompressed,optional"`
	IntervalBottom *float64          `db:"interval_bottom" json:"bottom,omitempty" parquet:"bottom,uncompressed,optional"`
	IntervalUnit   *string           `db:"interval_unit" json:"unit,omitempty" parquet:"unit,uncompressed,optional"`
	Keyword        []string          `json:"keyword,omitempty" parquet:"keyword,uncompressed,optional"`
	Barcode        *string           `json:"barcode,omitempty" parquet:"barcode,uncompressed,optional"`
	AltBarcode     []string          `json:"alt_barcode,omitempty" parquet:"alt_barcode,uncompressed,optional"`
	ContainerID    *int32            `json:"container_id,omitempty" parquet:"-"`
	ContainerPath  *string           `json:"container,omitempty" parquet:"container,uncompressed,optional"`
	Remark         *string           `json:"remark,omitempty" parquet:"-"`
	Geometries     []json.RawMessage `json:"geometries,omitempty" parquet:"-"`
	Latitude       *float64          `db:"latitude" json:"latitude,omitempty" parquet:"-"`
	Longitude      *float64          `db:"longitude" json:"longitude,omitempty" parquet:"-"`
	ProjectID      *int32            `json:"project_id,omitempty" parquet:"-"`
	Project        *string           `json:"project,omitempty" parquet:"project,uncompressed,optional"`
	CanPublish     *bool             `db:"can_publish" json:"can_publish,omitempty" parquet:"can_publish,uncompressed,optional"`
	Description    *string           `json:"description,omitempty" parquet:"description,uncompressed,optional"`
	Note           []string          `json:"note,omitempty" parquet:"-"`
	Issue          []string          `json:"issue,omitempty" parquet:"issue,uncompressed,optional"`
	Filename       []string          `json:"filename,omitempty" parquet:"-"`
	Interval       *struct {
		GTE *float64 `db:"gte" json:"gte,omitempty" parquet:"-"`
		LTE *float64 `db:"lte" json:"lte,omitempty" parquet:"-"`
	} `json:"interval,omitempty" parquet:"-"`
	Well []struct {
		ID        int32   `json:"id" parquet:"id,uncompressed"`
		Name      string  `json:"name,omitempty" parquet:"name,uncompressed,optional"`
		AltNames  string  `json:"alt_names,omitempty" parquet:"-"`
		Number    *string `json:"number,omitempty" parquet:"number,uncompressed,optional"`
		API       *string `json:"api,omitempty" parquet:"api,uncompressed,optional"`
		IsOnshore *bool   `json:"is_onshore,omitempty" parquet:"-"`
	} `json:"well,omitempty" parquet:"well,optional"`
	Outcrop []struct {
		ID     int32   `json:"id"  parquet:"id,uncompressed"`
		Name   string  `json:"name" parquet:"name,uncompressed"`
		Number *string `json:"number,omitempty" parquet:"number,uncompressed,optional"`
		Year   *int32  `json:"year,omitempty" parquet:"-"`
	} `json:"outcrop,omitempty" parquet:"outcrop,optional"`
	Borehole []struct {
		ID       int32  `json:"id" parquet:"id,uncompressed"`
		Name     string `json:"name,omitempty" parquet:"name,uncompressed,optional"`
		AltNames string `json:"alt_names,omitempty" parquet:"-"`
		Prospect struct {
			ID       int32  `json:"id" parquet:"id,uncompressed"`
			Name     string `json:"name,omitempty" parquet:"name,uncompressed,optional"`
			AltNames string `json:"alt_names,omitempty" parquet:"-"`
			ARDF     string `json:"ardf,omitempty" parquet:"-"`
		} `json:"prospect,omitempty" parquet:"prospect,optional"`
	} `json:"borehole,omitempty" parquet:"borehole,optional"`
	Shotline []struct {
		ID   int32    `json:"id" parquet:"id,uncompressed"`
		Name string   `json:"name" parquet:"name,uncompressed"`
		Year *int32   `json:"year,omitempty" parquet:"-"`
		Min  *float64 `json:"min,omitempty" parquet:"min,uncompressed,optional"`
		Max  *float64 `json:"max,omitempty" parquet:"max,uncompressed,optional"`
	} `json:"shotline,omitempty" parquet:"shotline,optional"`
	Publication []struct {
		ID          int32  `json:"id" parquet:"id,uncompressed"`
		Title       string `json:"title" parquet:"title,uncompressed"`
		Year        *int32 `json:"year,omitempty" parquet:"year,uncompressed,optional"`
		Description string `json:"description,omitempty" parquet:"-"`
		Number      string `json:"number,omitempty" parquet:"-"`
		Series      string `json:"series,omitempty" parquet:"-"`
	} `json:"publication,omitempty" parquet:"publication,optional"`
}

func (f *FlatInventory) StringID() string {
	return fmt.Sprintf("%d", f.ID)
}

func FlatInventoryFields(full bool) []string {
	if full {
		return []string{
			"ID",
			"Latitude",
			"Longitude",
			"Related",
			"Project",
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
		"Latitude",
		"Longitude",
		"Related",
		"Project",
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
		rel.WriteString(w.Name)
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
		if b.Prospect.Name != "" {
			rel.WriteString("Prospect: ")
			rel.WriteString(b.Prospect.Name)
			rel.WriteString("\n")
		}
		rel.WriteString("Borehole: ")
		rel.WriteString(b.Name)
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
			qfmt(f.Latitude),
			qfmt(f.Longitude),
			rel.String(),
			qfmt(f.Project),
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
			qfmt(f.Barcode),
			qfmt(f.ContainerPath),
			qfmt(f.Description),
		}
	}
	return []string{
		qfmt(f.ID),
		qfmt(f.Latitude),
		qfmt(f.Longitude),
		rel.String(),
		qfmt(f.Project),
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
