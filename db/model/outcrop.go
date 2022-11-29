package model

type Outcrop struct {
	ID            int32                    `json:"outcrop_id"`
	Name          string                   `json:"name"`
	IsOnshore     bool                     `json:"is_onshore"`
	Number        string                   `json:"outcrop_number"`
	Year          int16                    `json:"year"`
	Stash         map[string]interface{}   `json:"stash"`
	Keywords      []map[string]interface{} `json:"keywords"`
	GeoJSON       map[string]interface{}   `json:"geojson"`
	Quadrangles   []Quadrangle             `json:"quadrangles"`
	Notes         []Note                   `json:"notes"`
	URLs          []URL                    `json:"urls"`
	Organizations []Organization           `json:"organizations"`
	Files         []File                   `json:"files"`
}
