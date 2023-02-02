package model

type Shotline struct {
	ID             int32                  `json:"shotline_id"`
	Name           string                 `json:"name"`
	AltNames       string                 `json:"alt_names"`
	Year           int16                  `json:"year"`
	Remark         string                 `json:"remark"`
	ShotpointID    int32                  `json:"shotpoint_id"`
	Number         float64                `json:"shotpoint_number"`
	ShotpointMin   float64                `json:"shotpoint_min"`
	ShotpointMax   float64                `json:"shotpoint_max"`
	Stash          map[string]interface{} `json:"stash"`
	KeywordSummary []KeywordSummary       `json:"keywords"`
	GeoJSON        map[string]interface{} `json:"geojson"`
	Quadrangles    []Quadrangle           `json:"quadrangles"`
	Notes          []Note                 `json:"notes"`
	URLs           []URL                  `json:"urls"`
	Organizations  []Organization         `json:"organizations"`
	Files          []File                 `json:"files"`
}
