package model

type Prospect struct {
	ID         int32      `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	AltNames   *string    `db:"alt_names" json:"alt_names,omitempty"`
	ARDFNumber *string    `db:"ardf_number" json:"ardf_number,omitempty"`
	Boreholes  []Borehole `json:"boreholes,omitempty"`
	Files      []File     `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	KeywordSummary  []KeywordSummary `json:"keywords,omitempty"`
	GeoJSON         interface{}      `json:"geojson,omitempty"`
	MiningDistricts []MiningDistrict `db:"mining_districts" json:"mining_districts,omitempty"`
	Quadrangles     []Quadrangle     `json:"quadrangles,omitempty"`
}
