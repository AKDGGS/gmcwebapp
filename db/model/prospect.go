package model

type Prospect struct {
	ProspectID      int32                  `json:"prospect_id"`
	ProspectName    string                 `json:"prospect_name"`
	AltNames        string                 `json:"alt_names"`
	ARDFNumber      string                 `json:"ardf_number"`
	Files           []File                 `json:"files"`
	KeywordSummary  []KeywordSummary       `json:"keywords"`
	GeoJSON         map[string]interface{} `json:"geojson"`
	Quadrangles     []Quadrangle           `json:"quadrangles"`
	MiningDistricts []MiningDistricts      `json:"mining_district"`
	Boreholes       []Borehole             `json:"boreholes"`
}
