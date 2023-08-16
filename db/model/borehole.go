package model

type Borehole struct {
	ID                int32                  `json:"borehole_id,omitempty"`
	Name              string                 `json:"name,omitempty"`
	AltNames          string                 `json:"alt_name,omitempty"`
	Onshore           bool                   `json:"is_onshore"`
	CompletionDate    *FormattedDate         `json:"completion_date,omitempty"`
	MeasuredDepth     *float64               `json:"measured_depth,omitempty"`
	MeasuredDepthUnit string                 `json:"measured_depth_unit,omitempty"`
	Elevation         *float64               `json:"elevation,omitempty"`
	ElevationUnit     string                 `json:"elevation_unit,omitempty"`
	Stash             map[string]interface{} `json:"stash,omitempty"`
	EnteredDate       *FormattedDate         `json:"entered_date,omitempty"`
	ModifiedDate      *FormattedDate         `json:"modified_date,omitempty"`
	ModifiedUser      string                 `json:"modified_user,omitempty"`
	KeywordSummary    []KeywordSummary       `json:"keywords,omitempty"`
	GeoJSON           interface{}            `json:"geojson,omitempty"`
	MiningDistrict    MiningDistrict         `json:"mining_district,omitempty"`
	Quadrangles       []Quadrangle           `json:"quadrangles,omitempty"`
	Notes             []Note                 `json:"note,omitempty"`
	URLs              []URL                  `json:"url,omitempty"`
	Organizations     []Organization         `json:"organization,omitempty"`
	Files             []File                 `json:"file,omitempty"`
	Prospect          Prospect               `json:"prospect,omitempty"`
}
