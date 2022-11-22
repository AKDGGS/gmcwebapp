package model

type Outcrop struct {
	Name          string                 `json:"name"`
	ID            int32                  `json:"id"`
	IsOnshore     bool                   `json:"is_onshore"`
	OutcropNumber string                 `json:"outcrop_number"`
	Year          int16                  `json:"year"`
	Stash         map[string]interface{} `json:"stash"`
	Note          []Note                 `json:"note"`
	URL           []URL                  `json:"url"`
	Organization  []Organization         `json:"organization"`
	File          []File                 `json:"file"`
}
