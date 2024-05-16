package model

import (
	"encoding/json"
	"time"
)

type Well struct {
	ID               int32          `json:"id,omitempty"`
	Name             string         `json:"name,omitempty"`
	AltNames         *string        `db:"alt_names" json:"alt_names,omitempty"`
	Number           *string        `db:"number" json:"number,omitempty"`
	APINumber        *string        `db:"api_number" json:"api_number,omitempty"`
	Onshore          bool           `db:"onshore" json:"onshore"`
	Federal          bool           `db:"federal" json:"federal"`
	SpudDate         *time.Time     `db:"spud_date" json:"spud_date,omitempty"`
	CompletionDate   *time.Time     `db:"completion_date" json:"completion_date,omitempty"`
	MeasuredDepth    *float64       `db:"measured_depth" json:"measured_depth,omitempty"`
	VerticalDepth    *float64       `db:"vertical_depth" json:"vertical_depth,omitempty"`
	Elevation        *float64       `db:"elevation" json:"elevation,omitempty"`
	ElevationKB      *float64       `db:"elevation_kb" json:"elevation_kb,omitempty"`
	PermitStatus     *string        `db:"permit_status" json:"permit_status,omitempty"`
	PermitNumber     *int32         `db:"permit_number" json:"permit_number,omitempty"`
	CompletionStatus *string        `db:"completion_status" json:"completion_status,omitempty"`
	Unit             *string        `json:"unit,omitempty"`
	Stash            interface{}    `json:"stash,omitempty"`
	Notes            []Note         `json:"notes,omitempty"`
	URLs             []URL          `json:"urls,omitempty"`
	Organizations    []Organization `json:"organizations,omitempty"`
	Files            []File         `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
}

func (w *Well) MarshalJSON() ([]byte, error) {
	type Alias Well
	var spudDate string
	if w.SpudDate != nil {
		spudDate = w.SpudDate.Format("01-02-2006")
	}
	var completionDate string
	if w.CompletionDate != nil {
		completionDate = w.CompletionDate.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		SpudDate       string `json:"spud_date,omitempty"`
		CompletionDate string `json:"completion_date,omitempty"`
		*Alias
	}{
		SpudDate:       spudDate,
		CompletionDate: completionDate,
		Alias:          (*Alias)(w),
	})
}
