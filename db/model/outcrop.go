package model

import (
	"time"
)

type Outcrop struct {
	ID             int32                  `json:"outcrop_id"`
	Name           string                 `json:"name"`
	Number         string                 `json:"outcrop_number"`
	IsOnshore      bool                   `json:"is_onshore"`
	EnteredDate    *time.Time             `json:"entered_date"`
	ModifiedDate   *time.Time             `json:"modified_date"`
	ModifiedUser   string                 `json:"modified_user"`
	Year           int16                  `json:"year"`
	Stash          map[string]interface{} `json:"stash"`
	KeywordSummary []KeywordSummary       `json:"keywords"`
	GeoJSON        map[string]interface{} `json:"geojson"`
	Quadrangles    []Quadrangle           `json:"quadrangles"`
	Notes          []Note                 `json:"notes"`
	URLs           []URL                  `json:"urls"`
	Organizations  []Organization         `json:"organizations"`
	Files          []File                 `json:"files"`
}
