package model

type Outcrop struct {
	ID           int32                    `json:"outcrop_id"`
	Name         string                   `json:"name"`
	IsOnshore    bool                     `json:"is_onshore"`
	Number       string                   `json:"outcrop_number"`
	Year         int16                    `json:"year"`
	Stash        map[string]interface{}   `json:"stash"`
	Keywords     []map[string]interface{} `json:"keywords"`
	GeoJSON      map[string]interface{}   `json:"geojson"`
	Quadrangle   []Quadrangle             `json:"quadrangle"`
	Note         []Note                   `json:"note"`
	URL          []URL                    `json:"url"`
	Organization []Organization           `json:"organization"`
	File         []File                   `json:"file"`
}
