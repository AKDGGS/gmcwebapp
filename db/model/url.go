package model

type URL struct {
	ID          int32  `json:"url_id,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"url_type,omitempty"`
}
