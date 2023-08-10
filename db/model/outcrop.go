package model

import (
	"time"
)

type Outcrop struct {
	ID             int32                  `json:"outcrop_id,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Number         string                 `json:"outcrop_number,omitempty"`
	Onshore        bool                   `json:"is_onshore"`
	EnteredDate    *time.Time             `json:"entered_date,omitempty"`
	ModifiedDate   *time.Time             `json:"modified_date,omitempty"`
	ModifiedUser   string                 `json:"modified_user,omitempty"`
	Year           int16                  `json:"year,omitempty"`
	Stash          map[string]interface{} `json:"stash,omitempty"`
	KeywordSummary []KeywordSummary       `json:"keywords,omitempty"`
	GeoJSON        interface{}            `json:"geojson,omitempty"`
	Quadrangles    []Quadrangle           `json:"quadrangles,omitempty"`
	Notes          []Note                 `json:"notes,omitempty"`
	URLs           []URL                  `json:"urls,omitempty"`
	Organizations  []Organization         `json:"organizations,omitempty"`
	Files          []File                 `json:"files,omitempty"`
}
