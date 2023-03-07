package model

type Prospect struct {
	ID              int32             `json:"prospect_id"`
	Name            string            `json:"prospect_name"`
	AltNames        string            `json:"alt_names"`
	ARDFNumber      string            `json:"ardf_number"`
	Files           []File            `json:"files"`
	KeywordSummary  []KeywordSummary  `json:"keywords"`
	GeoJSON         interface{}       `json:"geojson"`
	Quadrangles     []Quadrangle      `json:"quadrangles"`
	MiningDistricts []MiningDistricts `json:"mining_district"`
	Boreholes       []Borehole        `json:"boreholes"`
}
