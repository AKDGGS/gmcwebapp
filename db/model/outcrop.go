package model

import (
	"encoding/json"
	"time"
)

type Outcrop struct {
	ID            int32          `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Number        *string        `db:"number" json:"number,omitempty"`
	Onshore       bool           `db:"onshore" json:"onshore"`
	EnteredDate   *time.Time     `db:"entered_date" json:"entered_date,omitempty"`
	ModifiedDate  *time.Time     `db:"modified_date" json:"modified_date,omitempty"`
	ModifiedUser  *string        `db:"modified_user" json:"modified_user,omitempty"`
	Year          *int16         `json:"year,omitempty"`
	Stash         interface{}    `json:"stash,omitempty"`
	Notes         []Note         `json:"notes,omitempty"`
	URLs          []URL          `json:"urls,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
	Files         []File         `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `db:"keywords" json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
}

func (o *Outcrop) MarshalJSON() ([]byte, error) {
	type Alias Outcrop
	var enteredDate string
	if o.EnteredDate != nil {
		enteredDate = o.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if o.ModifiedDate != nil {
		modifiedDate = o.ModifiedDate.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		EnteredDate  string `db:"entered_date" json:"entered_date,omitempty"`
		ModifiedDate string `db:"entered_date" json:"modified_date,omitempty"`
		*Alias
	}{
		EnteredDate:  enteredDate,
		ModifiedDate: modifiedDate,
		Alias:        (*Alias)(o),
	})
}
