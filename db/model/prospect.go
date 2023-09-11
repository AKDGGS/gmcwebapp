package model

type Prospect struct {
	ID         int32   `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	AltNames   *string `json:"alt_names,omitempty"`
	ARDFNumber *string `json:"ardf_number,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	Boreholes       []Borehole       `json:"boreholes,omitempty"`
	Files           []File           `json:"files,omitempty"`
	KeywordSummary  []KeywordSummary `json:"keywords,omitempty"`
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	MiningDistricts []MiningDistrict `json:"mining_districts,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
}
