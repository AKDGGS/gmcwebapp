package model

type URL struct {
	ID          int32  `json:"url_id"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Type        string `json:"url_type"`
}
