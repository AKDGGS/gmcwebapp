package model

import (
	"encoding/json"
	"time"
)

type Borehole struct {
	ID                int32                  `json:"id"`
	Name              *string                `json:"name"`
	AltNames          *string                `json:"alt_name,omitempty"`
	Onshore           bool                   `json:"is_onshore"`
	CompletionDate    time.Time              `json:"completion_date,omitempty"`
	MeasuredDepth     *float64               `json:"measured_depth,omitempty"`
	MeasuredDepthUnit *string                `json:"measured_depth_unit,omitempty"`
	Elevation         *float64               `json:"elevation,omitempty"`
	ElevationUnit     *string                `json:"elevation_unit,omitempty"`
	EnteredDate       time.Time              `json:"entered_date,omitempty"`
	ModifiedDate      time.Time              `json:"modified_date,omitempty"`
	ModifiedUser      *string                `json:"modified_user,omitempty"`
	Stash             map[string]interface{} `json:"stash,omitempty"`
	Notes             []Note                 `json:"notes,omitempty"`
	URLs              []URL                  `json:"urls,omitempty"`
	Organizations     []Organization         `json:"organizations,omitempty"`
	Files             []File                 `json:"files,omitempty"`
	Prospect          *Prospect              `json:"prospects,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	KeywordSummary  []KeywordSummary `json:"keywords,omitempty"`
	MiningDistricts []MiningDistrict `json:"mining_districts,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
}

func (b *Borehole) MarshalJSON() ([]byte, error) {
	type Alias Borehole
	var completionDate string
	if !b.CompletionDate.IsZero() {
		completionDate = b.CompletionDate.Format("01-02-2006")
	}
	var enteredDate string
	if !b.EnteredDate.IsZero() {
		enteredDate = b.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if !b.ModifiedDate.IsZero() {
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
