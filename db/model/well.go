package model

import (
	"encoding/json"
	"time"
)

type Well struct {
	ID               int32                  `json:"well_id"`
	Name             string                 `json:"name"`
	AltNames         string                 `json:"alt_name"`
	Number           string                 `json:"well_number"`
	APINumber        string                 `json:"api_number"`
	Onshore          bool                   `json:"is_onshore"`
	Federal          bool                   `json:"is_federal"`
	SpudDate         *time.Time             `json:"spud_date"`
	CompletionDate   *time.Time             `json:"completion_date"`
	MeasuredDepth    float64                `json:"measured_depth"`
	VerticalDepth    float64                `json:"vertical_depth"`
	Elevation        float64                `json:"elevation_depth"`
	ElevationKB      float64                `json:"elevation_kb"`
	PermitStatus     string                 `json:"permit_status"`
	PermitNumber     int32                  `json:"permit_number"`
	CompletionStatus string                 `json:"completion_status"`
	Unit             string                 `json:"unit"`
	Current          bool                   `json:"is_current"`
	OperatorName     string                 `json:"operator_name"`
	Remark           string                 `json:"remark"`
	OperatorType     string                 `json:"operator_type"`
	Stash            map[string]interface{} `json:"stash"`
	KeywordSummary   []KeywordSummary       `json:"keywords"`
	GeoJSON          map[string]interface{} `json:"geojson"`
	Quadrangles      []Quadrangle           `json:"quadrangles"`
	Notes            []Note                 `json:"note"`
	URLs             []URL                  `json:"url"`
	Organizations    []Organization         `json:"organization"`
	Files            []File                 `json:"file"`
}

func (w *Well) MarshalJSON() ([]byte, error) {
	type Alias Well
	return json.Marshal(&struct {
		CompletionDate string `json:"completion_date"`
		SpudDate       string `json:"spud_date"`
		*Alias
	}{
		CompletionDate: w.CompletionDate.Format("01-02-2006"),
		SpudDate:       w.SpudDate.Format("01-02-2006"),
		Alias:          (*Alias)(w),
	})
}
