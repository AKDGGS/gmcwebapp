package model

import (
	"encoding/json"
	"time"
)

type Borehole struct {
	ID                int32                  `json:"borehole_id"`
	Name              string                 `json:"name"`
	AltNames          string                 `json:"alt_name"`
	IsOnshore         bool                   `json:"is_onshore"`
	CompletionDate    *time.Time             `json:"completion_date"`
	MeasuredDepth     float64                `json:"measured_depth"`
	MeasuredDepthUnit string                 `json:"measured_depth_unit"`
	Elevation         float64                `json:"elevation"`
	ElevationUnit     string                 `json:"elevation_unit"`
	Stash             map[string]interface{} `json:"stash"`
	EnteredDate       *time.Time             `json:"entered_date"`
	ModifiedDate      *time.Time             `json:"modified_date"`
	ModifiedUser      string                 `json:"modified_user"`
	KeywordSummary    []KeywordSummary       `json:"keywords"`
	GeoJSON           map[string]interface{} `json:"geojson"`
	MiningDistricts   []MiningDistricts      `json:"mining_districts"`
	Quadrangles       []Quadrangle           `json:"quadrangles"`
	Notes             []Note                 `json:"note"`
	URLs              []URL                  `json:"url"`
	Organizations     []Organization         `json:"organization"`
	Files             []File                 `json:"file"`
	Prospect          Prospect               `json:"prospect"`
}

func (b *Borehole) MarshalJSON() ([]byte, error) {
	type Alias Borehole
	return json.Marshal(&struct {
		CompletionDate string `json:"completion_date"`
		*Alias
	}{
		CompletionDate: b.CompletionDate.Format("01-02-2006"),
		Alias:          (*Alias)(b),
	})
}
