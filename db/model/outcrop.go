package model

type Outcrop struct {
	ID           int32                  `json:"outcrop_id"`
	Name         string                 `json:"name"`
	IsOnshore    bool                   `json:"is_onshore"`
	Number       string                 `json:"outcrop_number"`
	Year         int16                  `json:"year"`
	Stash        map[string]interface{} `json:"stash"`
	Note         []Note                 `json:"note"`
	URL          []URL                  `json:"url"`
	Organization []Organization         `json:"organization"`
	File         []File                 `json:"file"`
}
