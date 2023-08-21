package model

type Prospect struct {
	ID              int32            `json:"prospect_id,omitempty"`
	Name            string           `json:"prospect_name,omitempty"`
	AltNames        string           `json:"alt_names,omitempty"`
	ARDFNumber      string           `json:"ardf_number,omitempty"`
	Files           []File           `json:"files,omitempty"`
	KeywordSummary  []KeywordSummary `json:"keywords,omitempty"`
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
	MiningDistricts []MiningDistrict `json:"mining_districts,omitempty"`
	Boreholes       []Borehole       `json:"boreholes,omitempty"`
}
