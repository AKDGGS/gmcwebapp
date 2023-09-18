package model

import (
	"encoding/json"
	"time"
)

type Outcrop struct {
	ID            int32                  `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Number        *string                `json:"outcrop_number,omitempty"`
	Onshore       bool                   `json:"is_onshore"`
	EnteredDate   time.Time              `json:"entered_date,omitempty"`
	ModifiedDate  time.Time              `json:"modified_date,omitempty"`
	ModifiedUser  *string                `json:"modified_user,omitempty"`
	Year          *int16                 `json:"year,omitempty"`
	Stash         map[string]interface{} `json:"stash,omitempty"`
	Notes         []Note                 `json:"notes,omitempty"`
	URLs          []URL                  `json:"urls,omitempty"`
	Organizations []Organization         `json:"organizations,omitempty"`
	Files         []File                 `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
}

func (o *Outcrop) MarshalJSON() ([]byte, error) {
	type Alias Outcrop
	var enteredDate string
	if !o.EnteredDate.IsZero() {
		enteredDate = o.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if !o.ModifiedDate.IsZero() {
		modifiedDate = o.ModifiedDate.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		EnteredDate  string `json:"entered_date,omitempty"`
		ModifiedDate string `json:"modified_date,omitempty"`
		*Alias
	}{
		EnteredDate:  enteredDate,
		ModifiedDate: modifiedDate,
		Alias:        (*Alias)(o),
	})
}
