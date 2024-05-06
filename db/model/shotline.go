package model

type Shotline struct {
	ID            int32                  `json:"id"`
	Name          string                 `json:"name"`
	AltNames      *string                `json:"alt_names,omitempty"`
	Year          *int16                 `json:"year,omitempty"`
	Remark        *string                `json:"remark,omitempty"`
	Stash         map[string]interface{} `json:"stash,omitempty"`
	Shotpoints    []Shotpoint            `json:"shotpoints,omitempty"`
	Notes         []Note                 `json:"notes,omitempty"`
	URLs          []URL                  `json:"urls,omitempty"`
	Organizations []Organization         `json:"organizations,omitempty"`
	Files         []File                 `json:"files,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	ShotpointMin   *float64         `json:"shotpoint_min,omitempty"`
	ShotpointMax   *float64         `json:"shotpoint_max,omitempty"`
	Numbers        *string          `json:"numbers,omitempty"`
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
}
