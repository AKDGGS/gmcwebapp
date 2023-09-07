package model

import (
	"encoding/json"
	"time"
)

type Well struct {
	ID               int32                  `json:"id,omitempty"`
	Name             *string                `json:"name,omitempty"`
	AltNames         *string                `json:"alt_names,omitempty"`
	Number           *string                `json:"well_number,omitempty"`
	APINumber        *string                `json:"api_number,omitempty"`
	Onshore          bool                   `json:"is_onshore"`
	Federal          bool                   `json:"is_federal"`
	SpudDate         *time.Time             `json:"spud_date,omitempty"`
	CompletionDate   *time.Time             `json:"completion_date,omitempty"`
	MeasuredDepth    *float64               `json:"measured_depth,omitempty"`
	VerticalDepth    *float64               `json:"vertical_depth,omitempty"`
	Elevation        *float64               `json:"elevation_depth,omitempty"`
	ElevationKB      *float64               `json:"elevation_kb,omitempty"`
	PermitStatus     *string                `json:"permit_status,omitempty"`
	PermitNumber     *int32                 `json:"permit_number,omitempty"`
	CompletionStatus *string                `json:"completion_status,omitempty"`
	Unit             *string                `json:"unit,omitempty"`
	Stash            map[string]interface{} `json:"stash,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
	Notes          []Note           `json:"notes,omitempty"`
	URLs           []URL            `json:"urls,omitempty"`
	Organizations  []Organization   `json:"organizations,omitempty"`
	Files          []File           `json:"files,omitempty"`
}

func (w *Well) MarshalJSON() ([]byte, error) {
	type Alias Well
	var spudDate string
	if w.SpudDate != nil && !w.SpudDate.IsZero() {
		spudDate = w.SpudDate.Format("01-02-2006")
	}
	var completionDate string
	if w.CompletionDate != nil && !w.CompletionDate.IsZero() {
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
