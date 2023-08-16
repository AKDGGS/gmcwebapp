package model

type Well struct {
	ID               int32                  `json:"well_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	AltNames         string                 `json:"alt_name,omitempty"`
	Number           string                 `json:"well_number,omitempty"`
	APINumber        string                 `json:"api_number,omitempty"`
	Onshore          bool                   `json:"is_onshore"`
	Federal          bool                   `json:"is_federal"`
	SpudDate         *FormattedDate         `json:"spud_date,omitempty"`
	CompletionDate   *FormattedDate         `json:"completion_date,omitempty"`
	MeasuredDepth    *float64               `json:"measured_depth,omitempty"`
	VerticalDepth    *float64               `json:"vertical_depth,omitempty"`
	Elevation        *float64               `json:"elevation_depth,omitempty"`
	ElevationKB      *float64               `json:"elevation_kb,omitempty"`
	PermitStatus     string                 `json:"permit_status,omitempty"`
	PermitNumber     *int32                 `json:"permit_number,omitempty"`
	CompletionStatus string                 `json:"completion_status,omitempty"`
	Unit             string                 `json:"unit,omitempty"`
	Stash            map[string]interface{} `json:"stash,omitempty"`
	KeywordSummary   []KeywordSummary       `json:"keywords,omitempty"`
	GeoJSON          interface{}            `json:"geojson,omitempty"`
	Quadrangles      []Quadrangle           `json:"quadrangles,omitempty"`
	Notes            []Note                 `json:"note,omitempty"`
	URLs             []URL                  `json:"url,omitempty"`
	Organizations    []Organization         `json:"organization,omitempty"`
	Files            []File                 `json:"file,omitempty"`
}
