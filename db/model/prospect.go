package model

type Prospect struct {
	ID         int32      `json:"prospect_id,omitempty"`
	Name       string     `json:"prospect_name,omitempty"`
	AltNames   string     `json:"alt_names,omitempty"`
	ARDFNumber string     `json:"ardf_number,omitempty"`
	Files      []File     `json:"files,omitempty"`
	Boreholes  []Borehole `json:"boreholes,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	KeywordSummary  []KeywordSummary `json:"keywords,omitempty"`
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	MiningDistricts []MiningDistrict `json:"mining_districts,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
}
