package model

type Shotline struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	AltNames      *string        `db:"alt_names" json:"alt_names,omitempty"`
	Year          *int16         `json:"year,omitempty"`
	Remark        *string        `json:"remark,omitempty"`
	Stash         interface{}    `json:"stash,omitempty"`
	Shotpoints    []Shotpoint    `json:"shotpoints,omitempty"`
	Notes         []Note         `json:"notes,omitempty"`
	URLs          []URL          `json:"urls,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
	Files         []File         `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	ShotpointMin   *float64         `db:"shotpoint_min" json:"shotpoint_min,omitempty"`
	ShotpointMax   *float64         `db:"shotpoint_max" json:"shotpoint_max,omitempty"`
	Numbers        *string          `json:"numbers,omitempty"`
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
}
