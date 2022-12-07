package model

import "time"

type Well struct {
	ID               int32                  `json:"well_id"`
	Name             string                 `json:"name"`
	AltName          string                 `json:"alt_name"`
	WellNumber       string                 `json:"well_number"`
	APINumber        string                 `json:"api_number"`
	IsOnshore        bool                   `json:"is_onshore"`
	IsFederal        bool                   `json:"is_federal"`
	PermitStatus     string                 `json:"permit_status"`
	CompletionStatus string                 `json:"completion_status"`
	SpudDate         time.Time              `json:"spud_date"`
	CompletionDate   time.Time              `json:"completion_date"`
	MeasuredDepth    float32                `json:"measured_depth"`
	VerticalDepth    float32                `json:"vertical_depth"`
	Elevation        float32                `json:"elevation_depth"`
	ElevationKb      float32                `json:"elevation_kb"`
	Stash            map[string]interface{} `json:"stash"`
	KeywordSummary   []KeywordSummary       `json:"keywords"`
	GeoJSON          map[string]interface{} `json:"geojson"`
	Quadrangles      []Quadrangle           `json:"quadrangles"`
	Notes            []Note                 `json:"note"`
	URLs             []URL                  `json:"url"`
	Organizations    []Organization         `json:"organization"`
	Files            []File                 `json:"file"`
}
