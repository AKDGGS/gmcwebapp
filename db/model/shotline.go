package model

type Shotline struct {
	ID          int32                  `json:"shotline_id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	AltNames    string                 `json:"alt_names,omitempty"`
	Year        *int16                 `json:"year,omitempty"`
	Remark      string                 `json:"remark,omitempty"`
	ShotpointID int32                  `json:"shotpoint_id,omitempty"`
	Number      *float64               `json:"shotpoint_number,omitempty"`
	Stash       map[string]interface{} `json:"stash,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON        interface{}      `json:"geojson,omitempty"`
	KeywordSummary []KeywordSummary `json:"keywords,omitempty"`
	Quadrangles    []Quadrangle     `json:"quadrangles,omitempty"`
	ShotpointMin   *float64         `json:"shotpoint_min,omitempty"`
	ShotpointMax   *float64         `json:"shotpoint_max,omitempty"`
	Shotpoints     []Shotpoint      `json:"shotpoints,omitempty"`
	Notes          []Note           `json:"notes,omitempty"`
	URLs           []URL            `json:"urls,omitempty"`
	Organizations  []Organization   `json:"organizations,omitempty"`
	Files          []File           `json:"files,omitempty"`
}
