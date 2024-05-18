package model

import (
	"encoding/json"
	"time"
)

type Borehole struct {
	ID                int32          `json:"id"`
	Name              *string        `json:"name"`
	AltNames          *string        `db:"alt_names" json:"alt_names,omitempty"`
	Onshore           bool           `db:"onshore" json:"onshore"`
	CompletionDate    *time.Time     `db:"completion_date" json:"completion_date,omitempty"`
	MeasuredDepth     *float64       `db:"measured_depth" json:"measured_depth,omitempty"`
	MeasuredDepthUnit *string        `db:"measured_depth_unit" json:"measured_depth_unit,omitempty"`
	Elevation         *float64       `json:"elevation,omitempty"`
	ElevationUnit     *string        `db:"elevation_unit" json:"elevation_unit,omitempty"`
	EnteredDate       *time.Time     `db:"entered_date" json:"entered_date,omitempty"`
	ModifiedDate      *time.Time     `db:"modified_date" json:"modified_date,omitempty"`
	ModifiedUser      *string        `db:"modified_user" json:"modified_user,omitempty"`
	Stash             interface{}    `json:"stash,omitempty"`
	Notes             []Note         `json:"notes,omitempty"`
	URLs              []URL          `json:"urls,omitempty"`
	Organizations     []Organization `json:"organizations,omitempty"`
	Files             []File         `json:"files,omitempty"`
	Prospect          *Prospect      `json:"prospect,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	KeywordSummary  []KeywordSummary `db:"keywords" json:"keywords,omitempty"`
	MiningDistricts []MiningDistrict `db:"mining_districts" json:"mining_districts,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
}

func (b *Borehole) MarshalJSON() ([]byte, error) {
	type Alias Borehole
	var completionDate string
	if b.CompletionDate != nil && !b.CompletionDate.IsZero() {
		completionDate = b.CompletionDate.Format("01-02-2006")
	}
	var enteredDate string
	if b.EnteredDate != nil && !b.EnteredDate.IsZero() {
		enteredDate = b.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if b.ModifiedDate != nil && !b.ModifiedDate.IsZero() {
		modifiedDate = b.ModifiedDate.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		CompletionDate string `json:"completion_date,omitempty"`
		EnteredDate    string `json:"entered_date,omitempty"`
		ModifiedDate   string `json:"modified_date,omitempty"`
		*Alias
	}{
		CompletionDate: completionDate,
		EnteredDate:    enteredDate,
		ModifiedDate:   modifiedDate,
		Alias:          (*Alias)(b),
	})
}
